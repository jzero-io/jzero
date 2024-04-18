package middlewares

import (
	"net/http"

	"google.golang.org/grpc/status"
)

func GrpcErrorHandler(err error) (int, any) {
	if st, ok := status.FromError(err); ok {
		return http.StatusOK, Body{
			Code: int(st.Code()),
			Msg:  st.Message(),
		}
	}

	code := http.StatusInternalServerError
	return http.StatusOK, Body{
		Code: code,
		Msg:  err.Error(),
	}
}
