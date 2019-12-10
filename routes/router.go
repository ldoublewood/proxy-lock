package routes

import (
	"github.com/gin-gonic/gin"
	"gin-gorm-example/controllers"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		proxyctl := new(controllers.ProxyController)
		v1.POST("/proxies", proxyctl.Import)
		v1.POST("/proxies/lock", proxyctl.Lock)
	}

	return router

}