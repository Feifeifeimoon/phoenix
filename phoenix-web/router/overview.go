package router

import (
	"app/common"
	"app/service"
	"github.com/gin-gonic/gin"
)

//整体运行状态
func sysInfo(c *gin.Context) {
	s, err := service.ServerStatus()
	if err != nil {
		c.JSON(200, common.ErrorResponse("获取服务运行状态失败"))
		return
	}

	ret := map[string]interface{}{
		"status":  s,
		"max_num": service.ServerMaxClient(),
		"cur_num": service.ServerCurClient(),
	}
	c.JSON(200, common.SuccessResponse(ret, 0))
	return
}

// 关闭，启动或者重启代理服务
func manageServer(c *gin.Context) {
	cmd, ok := c.GetQuery("command")
	if !ok {
		c.JSON(200, common.ErrorResponse("错误的参数"))
		return
	}
	if err := service.ManageServer(cmd); err != nil {
		c.JSON(200, common.ErrorResponse("管理服务失败"))
		return
	}
	c.JSON(200, common.SuccessResponse(nil, 0))
}
