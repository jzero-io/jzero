package hello

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jzero-io/jzero/daemon/internal/svc"
	"github.com/jzero-io/jzero/daemon/internal/types"
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
	time.Sleep(50 * time.Millisecond)
	resp = &types.Response{}
	resp.Message = fmt.Sprintf("Hello %s. I am Param", req.Name)
	return
}
