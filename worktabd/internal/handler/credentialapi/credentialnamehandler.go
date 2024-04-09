package credentialapi

import (
	"net/http"

	"github.com/jaronnie/worktab/worktabd/internal/logic/credentialapi"
	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CredentialNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := credentialapi.NewCredentialNameLogic(r.Context(), svcCtx)
		resp, err := l.CredentialName(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
