func new{{.upperStartCamelObject}}Model(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) *default{{.upperStartCamelObject}}Model {
	o := opts.DefaultApply(op...)
    var cachedConn sqlc.CachedConn
    if len(o.CacheConf) > 0 {
    	cachedConn = sqlc.NewConn(conn, o.CacheConf, o.CacheOpts...)
    }
    if o.CachedConn != nil {
    	cachedConn = *o.CachedConn
    }

    init{{.upperStartCamelObject}}Vars(o.Flavor)

	return &default{{.upperStartCamelObject}}Model{
		cachedConn: cachedConn,
		conn:       conn,
		flavor:     o.Flavor,
		table:      condition.QuoteWithFlavor(o.Flavor, "{{.data.Name.Source}}"),
	}
}

func (m *default{{.upperStartCamelObject}}Model) clone() *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		cachedConn: m.cachedConn,
		conn:       m.conn,
		table:      m.table,
		flavor:     m.flavor,
	}
}