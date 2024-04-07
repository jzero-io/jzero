package static

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaronnie/worktab/public"
)

func Static(r *gin.RouterGroup, f fs.FS) {
	staticHandler, err := fs.Sub(public.Public, "dist")
	if err != nil {
		log.Fatal("Unable to load static files: ", err)
	}
	r.StaticFS("/", http.FS(staticHandler))
}
