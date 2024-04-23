package machinelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jaronnie/jzero/daemon/internal/pb/machinepb"
	"github.com/jaronnie/jzero/daemon/internal/svc"
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
