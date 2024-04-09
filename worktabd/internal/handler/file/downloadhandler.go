package file

import (
	"net/http"

	"github.com/jaronnie/worktab/worktabd/internal/logic/file"
	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := file.NewDownloadLogic(r.Context(), svcCtx, w)
		err := l.Download(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
