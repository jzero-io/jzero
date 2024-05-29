package credentiallogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jzero-io/jzero/app/internal/pb/credentialpb"
	"github.com/jzero-io/jzero/app/internal/svc"
)

type CredentialVersion struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCredentialVersion(ctx context.Context, svcCtx *svc.ServiceContext) *CredentialVersion {
	return &CredentialVersion{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CredentialVersion) CredentialVersion(in *credentialpb.Empty) (*credentialpb.CredentialVersionResponse, error) {
	// todo: add your logic here and delete this line

	return &credentialpb.CredentialVersionResponse{}, nil
}
