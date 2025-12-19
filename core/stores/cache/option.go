package cache

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
)

// WithExpiry returns a func to customize an Options with given expiry.
func WithExpiry(expiry time.Duration) cache.Option {
	return func(o *cache.Options) {
		o.Expiry = expiry
	}
}

// WithNotFoundExpiry returns a func to customize an Options with given not found expiry.
func WithNotFoundExpiry(expiry time.Duration) cache.Option {
	return func(o *cache.Options) {
		o.NotFoundExpiry = expiry
	}
}
