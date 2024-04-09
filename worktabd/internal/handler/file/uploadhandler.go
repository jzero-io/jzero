package file

import (
	"net/http"

	"github.com/jaronnie/worktab/worktabd/internal/logic/file"
	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
