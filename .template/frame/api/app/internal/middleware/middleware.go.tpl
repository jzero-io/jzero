package middleware

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type Middleware struct {}

func NewMiddleware() Middleware {
	return Middleware{}
}

func Register(server *rest.Server) {
	httpx.SetOkHandler(ResponseMiddleware)
	httpx.SetErrorHandler(ErrorMiddleware)
	httpx.SetValidator(NewValidator())
}
