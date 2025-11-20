package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

type Body struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	return rw.body.Write(p)
}

func (rw *responseWriter) Body() []byte {
	return rw.body.Bytes()
}

func ErrorMiddleware(_ context.Context, err error) (int, any) {
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

func ResponseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logCtx := logx.ContextWithFields(r.Context(), logx.Field("path", r.URL.Path))

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		if strings.Contains(strings.ToLower(w.Header().Get("Content-Type")), "application/json") {
			var resp map[string]any
			err := json.Unmarshal(rw.Body(), &resp)
			if err != nil {
				logc.Errorf(logCtx, "Unmarshal resp error: %s\n", err.Error())
				return
			}

			if _, ok := resp["code"]; ok {
				httpx.OkJson(w, resp)
				return
			}

			wrappedResp := Body{
				Data:    resp,
				Code:    http.StatusOK,
				Message: "success",
			}
			httpx.OkJson(w, wrappedResp)
			return
		}

		_, err := w.Write(rw.Body())
		if err != nil {
			logc.Errorf(logCtx, "Write response error: %s\n", err.Error())
		}
	}
}
