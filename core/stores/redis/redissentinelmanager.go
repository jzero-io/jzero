package redis

import (
	"crypto/tls"
	"io"
	"runtime"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/syncx"
)

var (
	sentinelManager = syncx.NewResourceManager()

	// sentinelPoolSize is default pool size for sentinel type of redis.
	sentinelPoolSize = 5 * runtime.GOMAXPROCS(0)
)

func getSentinel(r *Redis) (*red.Client, error) {
	key := buildCacheKey(r.Addr, r.DB)
	val, err := sentinelManager.GetResource(key, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		// Parse sentinel addresses from Host field (comma-separated)
		sentinelAddrs := splitClusterAddrs(r.Addr)

		store := red.NewFailoverClient(&red.FailoverOptions{
			MasterName:       r.masterName,
			SentinelAddrs:    sentinelAddrs,
			Username:         r.User,
			Password:         r.Pass,
			DB:               r.DB,
			SentinelPassword: r.Pass, // Reuse Pass field for sentinel password
			MaxRetries:       maxRetries,
			MinIdleConns:     idleConns,
			TLSConfig:        tlsConfig,
		})

		hooks := append([]red.Hook{defaultDurationHook, breakerHook{
			brk: r.brk,
		}}, r.hooks...)
		for _, hook := range hooks {
			store.AddHook(hook)
		}

		connCollector.registerClient(&statGetter{
			clientType: SentinelType,
			key:        key,
			poolSize:   sentinelPoolSize,
			poolStats: func() *red.PoolStats {
				return store.PoolStats()
			},
		})

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
