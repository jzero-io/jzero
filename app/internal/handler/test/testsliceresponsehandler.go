package test

import (
	"net/http"

	"github.com/jzero-io/jzero/app/internal/logic/test"
	"github.com/jzero-io/jzero/app/internal/svc"
	"github.com/jzero-io/jzero/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TestSliceResponseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Empty
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := test.NewTestSliceResponseLogic(r.Context(), svcCtx)
		resp, err := l.TestSliceResponse(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
