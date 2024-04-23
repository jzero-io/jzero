package credentiallogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jaronnie/jzero/daemon/internal/model/credential"
	"github.com/jaronnie/jzero/daemon/internal/pb/credentialpb"
	"github.com/jaronnie/jzero/daemon/internal/svc"
)

type CreateCredentialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCredentialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCredentialLogic {
	return &CreateCredentialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCredentialLogic) CreateCredential(in *credentialpb.CreateCredentialRequest) (*credentialpb.CreateCredentialResponse, error) {
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
		Id:   id,
		Name: one.Name,
		Type: int32(one.Type),
	}, err
}
