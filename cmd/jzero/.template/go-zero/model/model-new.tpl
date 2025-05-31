func new{{.upperStartCamelObject}}Model(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) *default{{.upperStartCamelObject}}Model {
	o := opts.DefaultApply(op...)
    var cachedConn sqlc.CachedConn
    if len(o.CacheConf) > 0 {
    	cachedConn = sqlc.NewConn(conn, o.CacheConf, o.CacheOpts...)
    }
    if o.CachedConn != nil {
    	cachedConn = *o.CachedConn
    }

    initVars()

	return &default{{.upperStartCamelObject}}Model{
		cachedConn: cachedConn,
		conn: conn,
		table:      condition.Table({{.table}}),
	}
}