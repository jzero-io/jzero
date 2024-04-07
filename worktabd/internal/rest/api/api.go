package api

import "github.com/gin-gonic/gin"

func ApiRouter(rg *gin.RouterGroup) {
	rg.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "success")
	})
}
