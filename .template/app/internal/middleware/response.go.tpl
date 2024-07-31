package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
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

func ResponseHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logCtx := logx.ContextWithFields(r.Context(), logx.Field("path", r.URL.Path))

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		if strings.Contains(strings.ToLower(w.Header().Get("Content-Type")), "application/json") {
			var resp map[string]interface{}
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
