package router

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.Default()
	//中间层，提前做一些处理，比如日志
	//r.Use()

	//所有请求的路径从api
	api := r.Group("/api")
	{
		api.GET("/ping", ping)
		api.POST("/login", login)

		api.GET("/sys_info", sysInfo)
		api.POST("/manage_server", manageServer)
	}
	return r
}
