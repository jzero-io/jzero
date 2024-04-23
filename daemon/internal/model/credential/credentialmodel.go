package credential

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CredentialModel = (*customCredentialModel)(nil)

type (
	// CredentialModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCredentialModel.
	CredentialModel interface {
		credentialModel
		withSession(session sqlx.Session) CredentialModel
	}

	customCredentialModel struct {
		*defaultCredentialModel
	}
)

// NewCredentialModel returns a model for the database table.
func NewCredentialModel(conn sqlx.SqlConn) CredentialModel {
	return &customCredentialModel{
		defaultCredentialModel: newCredentialModel(conn),
	}
}

func (m *customCredentialModel) withSession(session sqlx.Session) CredentialModel {
	return NewCredentialModel(sqlx.NewSqlConnFromSession(session))
}
