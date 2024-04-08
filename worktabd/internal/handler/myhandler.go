package handler

import (
	"io/fs"
	"net/http"

	"github.com/jaronnie/worktab/public"
	"github.com/jaronnie/worktab/worktabd/internal/svc"
)

func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	}
}

func StaticFSHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	staticFS, _ := public.RootAssets()
	return dirhandler("/", staticFS)
}

func dirhandler(patern string, fs fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler := http.StripPrefix(patern, http.FileServer(http.FS(fs)))
		handler.ServeHTTP(w, req)
	}
}
