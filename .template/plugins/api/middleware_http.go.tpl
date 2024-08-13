package middleware

import (
    "net/http"
)

func {{.Name | FirstUpper}}Middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        next(writer, request)
    }
}