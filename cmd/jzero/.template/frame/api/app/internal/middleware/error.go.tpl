package middleware

import (
	"context"
	"net/http"
)

func ErrorMiddleware(_ context.Context, err error) (int, any) {
	return http.StatusOK, Body{
		Data:    nil,
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}
