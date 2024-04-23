package credentiallogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jaronnie/jzero/daemon/internal/pb/credentialpb"
	"github.com/jaronnie/jzero/daemon/internal/svc"
)

type CredentialVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCredentialVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CredentialVersionLogic {
	return &CredentialVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CredentialVersionLogic) CredentialVersion(in *credentialpb.Empty) (*credentialpb.CredentialVersionResponse, error) {
	// todo: add your logic here and delete this line

	return &credentialpb.CredentialVersionResponse{}, nil
}
