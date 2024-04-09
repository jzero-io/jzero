package machineapi

import (
	"net/http"

	"github.com/jaronnie/worktab/worktabd/internal/logic/machineapi"
	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func MachineNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := machineapi.NewMachineNameLogic(r.Context(), svcCtx)
		resp, err := l.MachineName(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
