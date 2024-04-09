// Code generated by goctl. DO NOT EDIT.
// Source: machine.proto

package server

import (
	"context"

	"github.com/jaronnie/worktab/worktabd/internal/logic/machinev2"
	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/pb/machinepb"
)

type Machinev2Server struct {
	svcCtx *svc.ServiceContext
	machinepb.UnimplementedMachinev2Server
}

func NewMachinev2Server(svcCtx *svc.ServiceContext) *Machinev2Server {
	return &Machinev2Server{
		svcCtx: svcCtx,
	}
}

func (s *Machinev2Server) MachineVersion(ctx context.Context, in *machinepb.Empty) (*machinepb.MachineVersionResponse, error) {
	l := machinev2logic.NewMachineVersionLogic(ctx, s.svcCtx)
	return l.MachineVersion(in)
}