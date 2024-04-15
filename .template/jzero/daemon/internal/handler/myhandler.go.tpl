package handler

import (
	"net/http"

	"{{ .Module }}/daemon/internal/svc"
)

func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("success"))
	}
}
