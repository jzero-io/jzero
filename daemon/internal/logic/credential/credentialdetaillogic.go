package credentiallogic

import (
	"context"

	"github.com/jaronnie/jzero/daemon/internal/pb/credentialpb"
	"github.com/jaronnie/jzero/daemon/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CredentialDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCredentialDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CredentialDetailLogic {
	return &CredentialDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CredentialDetailLogic) CredentialDetail(in *credentialpb.Int32Id) (*credentialpb.Credential, error) {
	// todo: add your logic here and delete this line

	return &credentialpb.Credential{}, nil
}
