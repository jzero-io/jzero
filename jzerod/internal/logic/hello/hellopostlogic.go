package hello

import (
	"context"
	"fmt"

	"github.com/jaronnie/jzero/jzerod/internal/svc"
	"github.com/jaronnie/jzero/jzerod/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelloPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloPostLogic {
	return &HelloPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloPostLogic) HelloPost(req *types.PostRequest) (resp *types.Response, err error) {
	resp = &types.Response{}
	resp.Message = fmt.Sprintf("Hello %s", req.Name)
	return
}
