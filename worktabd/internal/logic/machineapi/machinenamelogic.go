package machineapi

import (
	"context"

	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MachineNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMachineNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MachineNameLogic {
	return &MachineNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MachineNameLogic) MachineName(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
