package router

import (
	"app/common"
	"app/service"
	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	var s service.LoginService
	if err := c.ShouldBind(&s); err != nil {
		c.JSON(200, common.ErrorResponse(err.Error()))
		return
	}
	if err := s.Login(); err != nil {
		c.JSON(200, common.ErrorResponse("账号或密码错误"))
		return
	}
	c.JSON(200, common.SuccessResponse(nil, 0))
	return
}
