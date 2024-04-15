package hello

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/jaronnie/jzero/daemon/internal/types"
)

type HelloParamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloParamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloParamLogic {
	return &HelloParamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloParamLogic) HelloParam(req *types.ParamRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
