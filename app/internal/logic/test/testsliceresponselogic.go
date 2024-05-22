package testlogic

import (
	"context"

	"github.com/jzero-io/jzero/app/internal/svc"
	"github.com/jzero-io/jzero/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TestSliceResponseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestSliceResponseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestSliceResponseLogic {
	return &TestSliceResponseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestSliceResponseLogic) TestSliceResponse(req *types.Empty) (resp []types.TestResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
