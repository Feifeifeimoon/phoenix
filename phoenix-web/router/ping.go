package router

import (
	"app/common"
	"github.com/gin-gonic/gin"
)

func ping(ctx *gin.Context) {
	ctx.JSON(200, common.SuccessResponse(nil, 0))
}
