package credentiallogic

import (
	"context"
	"github.com/jaronnie/jzero/daemon/internal/model/credential"

	"github.com/jaronnie/jzero/daemon/internal/pb/credentialpb"
	"github.com/jaronnie/jzero/daemon/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CredentialListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCredentialListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CredentialListLogic {
	return &CredentialListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CredentialListLogic) CredentialList(in *credentialpb.CredentialListRequest) (*credentialpb.CredentialListResponse, error) {
	model := credential.NewCredentialModel(l.svcCtx.SqlConn)

	if in.GetPage() == 0 {
		in.Page = 1
	}

	if in.GetSize() == 0 {
		in.Size = 10
	}

	list, total, err := model.CredentialList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// trans []model.Credential to []credentialpb.Credential
	var pbList []*credentialpb.Credential

	for _, v := range list {
		pbList = append(pbList, &credentialpb.Credential{
			Id:   v.Id,
			Name: v.Name,
			Type: int32(v.Type),
		})
	}

	return &credentialpb.CredentialListResponse{
		List:  pbList,
		Total: int32(total),
	}, nil
}
