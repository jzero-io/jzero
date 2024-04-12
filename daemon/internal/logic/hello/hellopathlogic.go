package hello

import (
	"context"

	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/jaronnie/jzero/daemon/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelloPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloPathLogic {
	return &HelloPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloPathLogic) HelloPath(req *types.PathRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
