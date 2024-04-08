package routers

// routers TODO: Encapsulate routers for optimal use

import (
	"net/http"

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

	return routers
}
