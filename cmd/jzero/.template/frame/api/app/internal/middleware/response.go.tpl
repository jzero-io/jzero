package middleware

import (
	"context"
	"net/http"
)

type Body struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func ResponseMiddleware(_ context.Context, data any) any {
	return Body{
		Data:    data,
		Code:    http.StatusOK,
		Message: "success",
	}
}