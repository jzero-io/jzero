package credentiallogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jzero-io/jzero/app/internal/model/credential"
	"github.com/jzero-io/jzero/app/internal/pb/credentialpb"
	"github.com/jzero-io/jzero/app/internal/svc"
)

type CreateCredential struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCredential(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCredential {
	return &CreateCredential{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCredential) CreateCredential(in *credentialpb.CreateCredentialRequest) (*credentialpb.CreateCredentialResponse, error) {
	model := credential.NewCredentialModel(l.svcCtx.SqlConn)

	result, err := model.Insert(l.ctx, &credential.Credential{
		Name: in.Name,
		Type: int64(in.Type),
	})
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	one, err := model.FindOne(l.ctx, id)
	if err != nil {
		return nil, err
	}

	return &credentialpb.CreateCredentialResponse{
		Id:   int32(id),
		Name: one.Name,
		Type: int32(one.Type),
	}, err
}
