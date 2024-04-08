package routers

// routers TODO: Encapsulate routers for optimal use

import (
	"io/fs"
	"log"
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

	routers = append(routers, rest.Route{
		Method: "GET",
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/ui/", http.StatusMovedPermanently)
		},
	})

	// 添加静态文件服务路由
	staticHandler, err := fs.Sub(public.Public, "dist")
	if err != nil {
		log.Fatal("Unable to load static files: ", err)
	}

	routers = append(routers, rest.Route{
		Method: "GET",
		Path:   "/ui",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// 将 staticHandler 转换为 http.Handler
			// 使用 http.StripPrefix 去除路由前缀
			// 并将请求重定向到静态文件服务处理器
			http.StripPrefix("/ui", http.FileServer(http.FS(staticHandler))).ServeHTTP(w, r)
		},
	})

	return routers
}
