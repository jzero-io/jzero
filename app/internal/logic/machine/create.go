package machinelogic

import (
	"context"

	"github.com/jzero-io/jzero/app/internal/pb/machinepb"
	"github.com/jzero-io/jzero/app/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type Create struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreate(ctx context.Context, svcCtx *svc.ServiceContext) *Create {
	return &Create{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Create) Create(in *machinepb.Empty) (*machinepb.Empty, error) {
	// todo: add your logic here and delete this line

	return &machinepb.Empty{}, nil
}
