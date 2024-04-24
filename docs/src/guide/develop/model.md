---
title: model 数据库
icon: puzzle-piece
star: true
order: 5
category: 开发
tag:
  - Guide
---

jzero 推荐使用 go-zero sqlx 完成对数据库的 crud 操作.

jzero 数据库规范:

* sql 文件放在 daemon/desc/sql
* 生成的 model 放在 daemon/model

jzero gen 时会自动检测 daemon/desc/sql 下的 sql 文件, 并将生成的 model 放在 daemon/model 下

另外推荐使用 sqlbuilder 完成对 sql 的拼接. 如获取凭证列表以及支持过滤参数

```shell
go get github.com/huandu/go-sqlbuilder
```

```go
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
```