package hello

import (
	"context"

	"github.com/jaronnie/jzero/jzerod/internal/svc"
	"github.com/jaronnie/jzero/jzerod/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
