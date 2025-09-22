package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"

	"github.com/jzero-io/jzero/core/stores/redis"
)

const expiryDeviation = 0.05

type cacheNode struct {
	rds            *redis.Redis
	expiry         time.Duration
	notFoundExpiry time.Duration
	barrier        syncx.SingleFlight
	r              *rand.Rand
	lock           *sync.Mutex
	unstableExpiry mathx.Unstable
	stat           *Stat
	errNotFound    error
}

func (c cacheNode) Del(keys ...string) error {
	return c.Del(keys...)
}

func (c cacheNode) DelCtx(ctx context.Context, keys ...string) error {
	return c.DelCtx(ctx, keys...)
}

func (c cacheNode) Get(key string, val any) error {
	return c.GetCtx(context.Background(), key, val)
}

func (c cacheNode) GetCtx(ctx context.Context, key string, val any) error {
	return c.GetCtx(ctx, key, val)
}

func (c cacheNode) IsNotFound(err error) bool {
	return c.IsNotFound(err)
}

func (c cacheNode) Set(key string, val any) error {
	return c.SetCtx(context.Background(), key, val)
}

func (c cacheNode) SetCtx(ctx context.Context, key string, val any) error {
	return c.SetWithExpireCtx(ctx, key, val, c.expiry)
}

func (c cacheNode) SetWithExpire(key string, val any, expire time.Duration) error {
	return c.SetWithExpireCtx(context.Background(), key, val, expire)
}

func (c cacheNode) SetWithExpireCtx(ctx context.Context, key string, val any, expire time.Duration) error {
	return c.SetWithExpireCtx(ctx, key, val, expire)
}

func (c cacheNode) Take(val any, key string, query func(val any) error) error {
	return c.TakeCtx(context.Background(), val, key, query)
}

func (c cacheNode) TakeCtx(ctx context.Context, val any, key string, query func(val any) error) error {
	return c.TakeCtx(ctx, val, key, query)
}

func (c cacheNode) TakeWithExpire(val any, key string, query func(val any, expire time.Duration) error) error {
	return c.TakeWithExpireCtx(context.Background(), val, key, query)
}

func (c cacheNode) TakeWithExpireCtx(ctx context.Context, val any, key string, query func(val any, expire time.Duration) error) error {
	return c.TakeWithExpireCtx(ctx, val, key, query)
}

// NewNode returns a cacheNode.
// rds is the underlying redis node or cluster.
// barrier is the barrier that maybe shared with other cache nodes on cache cluster.
// st is used to stat the cache.
// errNotFound defines the error that returned on cache not found.
// opts are the options that customize the cacheNode.
func NewNode(rds *redis.Redis, barrier syncx.SingleFlight, st *Stat,
	errNotFound error, opts ...Option,
) cache.Cache {
	o := newOptions(opts...)
	return cacheNode{
		rds:            rds,
		expiry:         o.Expiry,
		notFoundExpiry: o.NotFoundExpiry,
		barrier:        barrier,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           st,
		errNotFound:    errNotFound,
	}
}
