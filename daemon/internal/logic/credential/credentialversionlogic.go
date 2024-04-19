package credentiallogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/jaronnie/jzero/daemon/pb/credentialpb"
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
	time.Sleep(50 * time.Millisecond)
	return &credentialpb.CredentialVersionResponse{
		Version: "v1.0",
	}, nil
}
