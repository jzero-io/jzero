package credentialapi

import (
	"context"

	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CredentialNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCredentialNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CredentialNameLogic {
	return &CredentialNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CredentialNameLogic) CredentialName(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
