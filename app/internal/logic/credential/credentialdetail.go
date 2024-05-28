package credentiallogic

import (
	"context"

	"github.com/jzero-io/jzero/app/internal/pb/credentialpb"
	"github.com/jzero-io/jzero/app/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CredentialDetail struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCredentialDetail(ctx context.Context, svcCtx *svc.ServiceContext) *CredentialDetail {
	return &CredentialDetail{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CredentialDetail) CredentialDetail(in *credentialpb.Int32Id) (*credentialpb.Credential, error) {
	// todo: add your logic here and delete this line

	return &credentialpb.Credential{}, nil
}
