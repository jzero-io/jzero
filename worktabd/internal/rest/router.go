package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jaronnie/worktab/public"
	"github.com/jaronnie/worktab/worktabd/internal/rest/api"
	"github.com/jaronnie/worktab/worktabd/internal/rest/static"
)

func Router(e *gin.Engine) *gin.Engine {
	// redirect 到 /@manage
	e.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/@manage")
	})

	ui := e.Group("/@manage")
	static.Static(ui, public.Public)

	apiV1 := e.Group("/api/v1")
	api.ApiRouter(apiV1)

	return e
}
