func new{{.upperStartCamelObject}}Model(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) *default{{.upperStartCamelObject}}Model {
	{{if .withCache}}o := opts.DefaultApply(op...){{end}}
	return &default{{.upperStartCamelObject}}Model{
		{{if .withCache}}cachedConn: sqlc.NewConn(conn, o.CacheConf, o.CacheOpts...),{{end}}
		conn: conn,
		table:      {{.table}},
	}
}

