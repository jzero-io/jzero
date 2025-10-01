package modelx

import (
	"github.com/eddieowens/opts"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type ModelOpts struct {
	CachedConn *sqlc.CachedConn
	CacheConf  cache.CacheConf
	CacheOpts  []cache.Option
}

func (opts ModelOpts) DefaultOptions() ModelOpts {
	return ModelOpts{}
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
