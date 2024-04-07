package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jaronnie/worktab/public"
	"github.com/jaronnie/worktab/worktabd/internal/rest/api"
	"github.com/jaronnie/worktab/worktabd/internal/rest/static"
)

func Router(g *gin.Engine) *gin.Engine {
	// redirect 到 /@manage
	g.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/@manage")
	})

	ui := g.Group("/@manage")
	static.Static(ui, public.Public)

	apiV1 := g.Group("/api/v1")
	api.ApiRouter(apiV1)

	return g
}
