package handler

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest"

	"github.com/jzero-io/jzero/daemon/internal/svc"
	"github.com/jzero-io/jzero/embeded"
)

func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("success"))
	}
}

func registerStaticEmbedHandler(server *rest.Server, serverCtx *svc.ServiceContext) {
	// related: https://blog.csdn.net/keytounix/article/details/108424389
	dirLevel := []string{":1", ":2", ":3", ":4", ":5", ":6", ":7", ":8"}
	pattern := "/"
	staticFS, _ := embeded.RootWeb()
	for i := 1; i < len(dirLevel); i++ {
		path := "/" + strings.Join(dirLevel[:i], "/")
		server.AddRoute(
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: dirHandler(pattern, staticFS),
			})
	}
}

func dirHandler(pattern string, fs fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler := http.StripPrefix(pattern, http.FileServer(http.FS(fs)))
		handler.ServeHTTP(w, req)
	}
}
