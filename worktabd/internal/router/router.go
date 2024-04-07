package router

import "github.com/gin-gonic/gin"

func Load(g *gin.Engine) *gin.Engine {
	g.GET("/api/v2/hello", func(ctx *gin.Context) {
		ctx.String(200, "hello world")
	})
	return g
}
