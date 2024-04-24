package credential

import (
	"context"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jaronnie/jzero/daemon/internal/pb/credentialpb"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CredentialModel = (*customCredentialModel)(nil)

type (
	// CredentialModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCredentialModel.
	CredentialModel interface {
		credentialModel
		withSession(session sqlx.Session) CredentialModel

		CredentialList(ctx context.Context, options *credentialpb.CredentialListRequest) ([]*Credential, int64, error)
	}

	customCredentialModel struct {
		*defaultCredentialModel
	}
)

func (m *customCredentialModel) CredentialList(ctx context.Context, options *credentialpb.CredentialListRequest) ([]*Credential, int64, error) {
	sb := sqlbuilder.Select("*").From(m.table)
	countsb := sqlbuilder.Select("count(*)").From(m.table)

	if options.GetName() != "" {
		sb.Where(sb.Like("name", options.GetName()))
		countsb.Where(sb.Like("name", options.GetName()))
	}

	sb.Limit(int(options.GetSize())).Offset(int((options.GetPage() - 1) * options.GetSize()))

	var credentials []*Credential

	err := m.conn.QueryRowsCtx(ctx, &credentials, sb.String())
	if err != nil {
		return nil, 0, err
	}

	// get total
	var total int64
	err = m.conn.QueryRowCtx(ctx, &total, countsb.String())
	if err != nil {
		return nil, 0, err
	}

	return credentials, total, nil
}

// NewCredentialModel returns a model for the database table.
func NewCredentialModel(conn sqlx.SqlConn) CredentialModel {
	return &customCredentialModel{
		defaultCredentialModel: newCredentialModel(conn),
	}
}

func (m *customCredentialModel) withSession(session sqlx.Session) CredentialModel {
	return NewCredentialModel(sqlx.NewSqlConnFromSession(session))
}
