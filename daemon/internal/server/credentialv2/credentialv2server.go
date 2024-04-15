// Code generated by goctl. DO NOT EDIT.
// Source: credential.proto

package server

import (
	"context"

	credentialv2logic "github.com/jaronnie/jzero/daemon/internal/logic/credentialv2"
	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/jaronnie/jzero/daemon/pb/credentialpb"
)

type Credentialv2Server struct {
	svcCtx *svc.ServiceContext
	credentialpb.UnimplementedCredentialv2Server
}

func NewCredentialv2Server(svcCtx *svc.ServiceContext) *Credentialv2Server {
	return &Credentialv2Server{
		svcCtx: svcCtx,
	}
}

func (s *Credentialv2Server) CredentialVersion(ctx context.Context, in *credentialpb.Empty) (*credentialpb.CredentialVersionResponse, error) {
	l := credentialv2logic.NewCredentialVersionLogic(ctx, s.svcCtx)
	return l.CredentialVersion(in)
}