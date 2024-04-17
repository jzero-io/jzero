package hello

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/jaronnie/jzero/daemon/internal/types"
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
	resp = &types.Response{}
	resp.Message = fmt.Sprintf("Hello %s. I am Path", req.Name)
	return
}
