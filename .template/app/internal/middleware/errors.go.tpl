package middleware

import (
	"net/http"

	"google.golang.org/grpc/status"
)

func ErrorHandler(err error) (int, any) {
	code := http.StatusInternalServerError
	message := err.Error()

	// from grpc error
	if st, ok := status.FromError(err); ok {
		code = int(st.Code())
		message = st.Message()
	}

	return http.StatusOK, Body{
		Data:    nil,
		Code:    code,
		Message: message,
	}
}
