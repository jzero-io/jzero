package test

import (
	"context"

	"github.com/jzero-io/jzero/app/internal/svc"
	"github.com/jzero-io/jzero/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TestSliceResponse2 struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestSliceResponse2(ctx context.Context, svcCtx *svc.ServiceContext) *TestSliceResponse2 {
	return &TestSliceResponse2{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestSliceResponse2) TestSliceResponse2(req *types.Empty) (resp []types.TestResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
