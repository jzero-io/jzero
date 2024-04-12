package machinev2logic

import (
	"context"

	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/jaronnie/jzero/daemon/pb/machinepb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MachineVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMachineVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MachineVersionLogic {
	return &MachineVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MachineVersionLogic) MachineVersion(in *machinepb.Empty) (*machinepb.MachineVersionResponse, error) {
	// todo: add your logic here and delete this line

	return &machinepb.MachineVersionResponse{}, nil
}
