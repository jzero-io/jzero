package cache

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Cache interface {
	cache.Cache

	// SetNoExpireCtx Because zero cache set ctx has default expire, so jzero add this method
	SetNoExpireCtx(ctx context.Context, key string, val any) error

	// GetPrefixKeysCtx get prefix key, give prefix key return all matched key
	GetPrefixKeysCtx(ctx context.Context, prefix string) ([]string, error)

	// ExpireCtx set key expire
	ExpireCtx(ctx context.Context, key string, expire time.Duration) error
}
