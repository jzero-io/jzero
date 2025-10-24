package modelx

import (
	"github.com/eddieowens/opts"
	"github.com/huandu/go-sqlbuilder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type ModelOpts struct {
	CachedConn *sqlc.CachedConn
	CacheConf  cache.CacheConf
	CacheOpts  []cache.Option
	Flavor     sqlbuilder.Flavor
}

func (opts ModelOpts) DefaultOptions() ModelOpts {
	return ModelOpts{
		Flavor: sqlbuilder.DefaultFlavor,
	}
}

func WithCachedConn(cachedConn sqlc.CachedConn) opts.Opt[ModelOpts] {
	return func(o *ModelOpts) {
		o.CachedConn = &cachedConn
	}
}

func WithCacheConf(cacheConf cache.CacheConf) opts.Opt[ModelOpts] {
	return func(o *ModelOpts) {
		o.CacheConf = cacheConf
	}
}

func WithCacheOpts(cacheOpts ...cache.Option) opts.Opt[ModelOpts] {
	return func(o *ModelOpts) {
		o.CacheOpts = cacheOpts
	}
}

func WithFlavor(flavor sqlbuilder.Flavor) opts.Opt[ModelOpts] {
	return func(o *ModelOpts) {
		o.Flavor = flavor
	}
}
