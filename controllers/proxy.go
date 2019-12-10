package controllers

import (
	"gin-gorm-example/database"
	"gin-gorm-example/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)
const (

)

const (
	DefaultDurationMinute = int64(60)
)
type ProxyController struct {
	Basic
}
func (a *ProxyController) Create(c *gin.Context, proxy models.Proxy) (errcode int, err error) {
	// -1 is default json value, so means not provided
	if proxy.DurationMinute == -1 {
		proxy.DurationMinute = DefaultDurationMinute
	}
	dur := int64(time.Minute) * proxy.DurationMinute
	proxy.ExpiredAt = time.Now().Add(time.Duration(dur))

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
	a.JsonSuccess(c, http.StatusCreated, gin.H{"data": *proxy})
}

func (a *ProxyController) doImport(c *gin.Context) (errcode int, err error) {
	var request ImportRequest
	if err = c.ShouldBind(&request); err == nil {
		for _, p := range request.Proxys {
			if errcode, err = a.Create(c, p); err != nil {
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

func (a *ProxyController) doLock(c *gin.Context) (errcode int, err error, proxy *models.Proxy) {
	var request TakeRequest
	if err = c.ShouldBind(&request); err == nil {
		var proxy models.Proxy
		if err := database.DB.Where(&models.Proxy{
			Status:         "free",
		}).Take(&proxy).Error; err != nil {
			log.Printf("no avaiable proxy!! please contact admin")
			return http.StatusNotFound, err, nil
		}
		proxy.Status = "locked"
		proxy.Owner = request.Owner

		if err = database.DB.Model(&models.Proxy{}).Update(&proxy).Error; err != nil {
			log.Printf("lock the record fail!! please contact admin")
			return http.StatusNotFound, err, nil
		}
		//a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "OK"})
		return http.StatusOK, nil, &proxy
	} else {
		//a.JsonFail(c, http.StatusBadRequest, err.Error())
		return http.StatusBadRequest, err, nil
	}
}


type ImportRequest struct {
	Proxys []models.Proxy `form:"proxys" json:"proxys" binding:"required"`
	Tag		string `form:"tag" json:"tag" binding:"required"`
}


type TakeRequest struct {
	Owner 	string `form:"owner" json:"owner" binding:"required"`
}




