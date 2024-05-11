package middlewares

import (
	"net/http"

	"google.golang.org/grpc/status"
)

func ErrorHandler(err error) (int, any) {
	if st, ok := status.FromError(err); ok {
		return int(st.Code()), err
	}

	code := http.StatusInternalServerError
	return code, err
}
