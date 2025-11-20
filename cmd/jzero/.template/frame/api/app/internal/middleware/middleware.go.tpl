package middleware

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func Register(server *rest.Server) {
	httpx.SetOkHandler(NewOkMiddleware().Handle)
	httpx.SetErrorHandlerCtx(NewErrorMiddleware().Handle)
	httpx.SetValidator(NewValidator())

	// add custom middleware
	// server.Use(func(next http .HandlerFunc) http.HandlerFunc {
	//	return func(w http.ResponseWriter, r *http.Request) {
	//		next.ServeHTTP(w, r)
	//	}
	// })
}
