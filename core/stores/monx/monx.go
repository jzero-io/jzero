package monx

import (
	"github.com/eddieowens/opts"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type MonOpts struct {
	CacheConf cache.CacheConf
	CacheOpts []cache.Option
}

func (opts MonOpts) DefaultOptions() MonOpts {
	return MonOpts{}
}

func WithCacheConf(cacheConf cache.CacheConf) opts.Opt[MonOpts] {
	return func(o *MonOpts) {
		o.CacheConf = cacheConf
	}
}

func WithCacheOpts(cacheOpts ...cache.Option) opts.Opt[MonOpts] {
	return func(o *MonOpts) {
		o.CacheOpts = cacheOpts
	}
}
