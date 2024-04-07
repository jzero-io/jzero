package logic

import (
	"context"

	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/worktabdpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServerVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewServerVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServerVersionLogic {
	return &ServerVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ServerVersionLogic) ServerVersion(in *worktabdpb.Empty) (*worktabdpb.ServerVersionResponse, error) {
	// todo: add your logic here and delete this line

	return &worktabdpb.ServerVersionResponse{
		Version: "v1",
	}, nil
}
