package middleware

import (
	"app/common"
	"github.com/gin-gonic/gin"
)

/*
	jwt验证
*/

func JwtVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(401, common.ErrorResponse("未携带token"))
			c.Abort()
			return
		}
		_, err := common.VerifyToken(token)
		if err != nil {
			if err == common.TokenExpired {
				c.JSON(401, common.ErrorResponse(err.Error()))
				c.Abort()
			} else {
				c.JSON(401, common.ErrorResponse(err.Error()))
				c.Abort()
			}
			return
		}
		//var user dao.User
		//if err := dao.DB.Where("id = ?", id).First(&user).Error; err != nil {
		//	c.JSON(200, util.ErrorResponse(err.Error()))
		//	c.Abort()
		//	return
		//}
		//c.Set("user", &user)
		c.Next()
	}
}
