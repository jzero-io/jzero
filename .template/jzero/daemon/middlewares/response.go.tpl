package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type Body struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
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

func WrapResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logCtx := logx.ContextWithFields(r.Context(), logx.Field("path", r.URL.Path))

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		if rw.statusCode != http.StatusOK {
			return
		}

		var resp map[string]interface{}
		err := json.Unmarshal(rw.Body(), &resp)
		if err != nil {
			_, err := w.Write(rw.Body())
			if err != nil {
				logc.Errorf(logCtx, "Write response error: %s\n", err.Error())
			}
			return
		}

		if _, ok := resp["code"]; ok {
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(rw.Body())
			if err != nil {
				logc.Errorf(logCtx, "Write response error: %s\n", err.Error())
			}
			return
		}

		wrappedResp := map[string]interface{}{
			"code":    http.StatusOK,
			"message": "success",
			"data":    resp,
		}
		httpx.OkJson(w, wrappedResp)
	}
}
