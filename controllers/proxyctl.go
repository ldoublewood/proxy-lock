package controllers

import (
	"fmt"
	"gin-gorm-example/config"
	"gin-gorm-example/database"
	"gin-gorm-example/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultDurationMinute = int64(30)
)
type ProxyController struct {
	Basic
}
func (a *ProxyController) Create(c *gin.Context, proxy models.Proxy, now time.Time) (errcode int, err error) {
	// 0 is default json value, so means not provided
	if proxy.DurationMinute == 0 {
		proxy.DurationMinute = DefaultDurationMinute
	}
	dur := int64(time.Minute) * proxy.DurationMinute
	proxy.ExpiredAt = now.Add(time.Duration(dur))

	if err := database.DB.Create(&proxy).Error; err != nil {
//		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return http.StatusServiceUnavailable, err
	}
	return http.StatusOK, nil
}

func (a *ProxyController) Import(c *gin.Context) {
	if errcode, err := a.doImport(c); err != nil {
		a.JsonFail(c, errcode, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "OK"})
}

func (a *ProxyController) Lock(c *gin.Context) {
	var proxy *models.Proxy
	var errcode int
	var err error
	if errcode, err, proxy = a.doLock(c); err != nil {
		a.JsonFail(c, errcode, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"data": *proxy})
}

func (a *ProxyController) doImport(c *gin.Context) (errcode int, err error) {
	var request ImportRequest
	var now = time.Now()
	if err = c.ShouldBind(&request); err == nil {
		for _, p := range request.Proxies {
			if errcode, err = a.Create(c, p, now); err != nil {
				return errcode, err
			}
		}
		//a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "OK"})
		return http.StatusOK, nil
	} else {
		//a.JsonFail(c, http.StatusBadRequest, err.Error())
		return http.StatusBadRequest, err
	}
}

func scanExpiredRecord(now time.Time) (errcode int, err error) {

	if err = database.DB.Model(&models.Proxy{}).Where("status = 'free' AND expired_at <= ?", now).Update(&models.Proxy{
		Status:         "expired",
	}).Error; err != nil {
		log.Printf("scan expire record fail!! please contact admin")
		return http.StatusServiceUnavailable, err
	}
	return http.StatusOK, nil

}

func getTestedProxy(owner string, destUrl string) (errcode int, err error, outproxy *models.Proxy) {
	var proxy models.Proxy
	for {
		if err = database.DB.Where(&models.Proxy{
			Status: "free",
		}).Take(&proxy).Error; err != nil {
			log.Printf("no avaiable proxy!! please contact admin")
			return http.StatusNotFound, err, nil
		}

		proxy.Status = "locked"
		proxy.Owner = owner

		if err = database.DB.Model(&models.Proxy{}).Update(&proxy).Error; err != nil {
			log.Printf("lock the record fail!! please contact admin")
			return http.StatusNotFound, err, nil
		}

		if errcode, err = doTest(destUrl, &proxy); err != nil {
			log.Printf("test request fail, so mask the record and try next record: %s", err.Error())
			proxy.Status = "masked"
			proxy.Owner = owner
			if err = database.DB.Model(&models.Proxy{}).Update(&proxy).Error; err != nil {
				log.Printf("mask the record fail")
				return http.StatusNotFound, err, nil
			}
			continue
		}
		log.Printf("found tested proxy: %s:%d", proxy.Host, proxy.Port)
		break
	}


	return http.StatusOK, nil, &proxy

}

func doTest(destUrl string, proxy *models.Proxy) (errcode int, err error) {
	// return OK if destUrl is empty
	if destUrl == "" {
		return http.StatusOK, nil
	}
	urli := url.URL{}
	strProxy := fmt.Sprintf("http://%s:%d", proxy.Host, proxy.Port)
	urlproxy, _ := urli.Parse(strProxy)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		},
	}
	var req *http.Request
	req, err = http.NewRequest("GET", destUrl, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}
	var response *http.Response
	response, err = client.Do(req)
	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	//ignore the response because this is only test
	return http.StatusOK, nil
}
func (a *ProxyController) doLock(c *gin.Context) (errcode int, err error, outproxy *models.Proxy) {
	var request TakeRequest

	if err = c.ShouldBind(&request); err != nil {
		return http.StatusBadRequest, err, nil
	}
	var now = time.Now()
	if errcode, err = scanExpiredRecord(now); err != nil {
		return
	}
	return getTestedProxy(request.Owner, config.Get().TestDestUrl)
}


type ImportRequest struct {
	Proxies []models.Proxy `form:"proxies" json:"proxies" binding:"required"`
	Tag		string `form:"tag" json:"tag" binding:"required"`
}


type TakeRequest struct {
	Owner 	string `form:"owner" json:"owner" binding:"required"`
}




