---
title: model 数据库教程
icon: icon-park-twotone:database-code
star: true
order: 5
category: 开发
tag:
  - Guide
---

jzero 推荐使用 go-zero sqlx 完成对数据库的 crud 操作.

jzero 数据库规范:

* sql 文件放在 desc/sql
* 生成的 model 放在 internal/model

jzero gen 时会自动检测 desc/sql 下的 sql 文件, 并将生成的 model 放在 internal/model 下

另外推荐使用 sql-builder 完成对 sql 的拼接. 如获取凭证列表以及支持过滤参数

## 将默认生成的 crud 代码替换为使用 sql-builder 的方式

```shell
git clone https://github.com/jzero-io/sqlbuilder-zero
cp -r sqlbuilder-zero/templates/go-zero/model ~/.jzero/$Version/go-zero
```

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