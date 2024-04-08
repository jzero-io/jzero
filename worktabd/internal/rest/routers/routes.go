package routers

// routers TODO: Encapsulate routers for optimal use

import (
	"io/fs"
	"net/http"

	"github.com/jaronnie/worktab/public"
	"github.com/zeromicro/go-zero/rest"
)

func SetRoutes() []rest.Route {
	var routers []rest.Route

	routers = append(routers, rest.Route{
		Method: "GET",
		Path:   "/api/v1.0/health",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("success"))
		},
	})

	staticFS, _ := public.RootAssets()

	routers = append(routers, rest.Route{
		Method:  "GET",
		Path:    "/",
		Handler: dirhandler("/", staticFS),
	})
	return routers
}

func dirhandler(patern string, fs fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler := http.StripPrefix(patern, http.FileServer(http.FS(fs)))
		handler.ServeHTTP(w, req)
	}
}
