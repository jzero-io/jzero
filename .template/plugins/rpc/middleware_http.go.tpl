package middleware

import (
    "fmt"
    "net/http"

    "github.com/zeromicro/go-zero/core/logx"
)

func {{.Name | FirstUpper}}Middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        if MatchRoute("{{.Name}}", fmt.Sprintf("%s:%s",request.Method,request.URL.Path)) {
            // do something before middleware
            logx.WithContext(request.Context()).Info("enter {{.Name}} before middleware")
        }
        next.ServeHTTP(writer, request)
        if MatchRoute("{{.Name}}", fmt.Sprintf("%s:%s",request.Method,request.URL.Path)) {
            // do something after middleware
            logx.WithContext(request.Context()).Info("enter {{.Name}} after middleware")
        }
    }
}