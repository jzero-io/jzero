package cache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestRedisNode(t *testing.T) {
	r, err := miniredis.Run()
	assert.NoError(t, err)
	defer r.Close()

	cache := NewRedisNode(redis.New(r.Addr()), errors.New("not found"))

	err = cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:abc", "abc")
	assert.NoError(t, err)

	err = cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:def", "def")
	assert.NoError(t, err)

	err = cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:ghi", "ghi")
	assert.NoError(t, err)

	keys, err := cache.GetPrefixKeysCtx(context.Background(), "JWT_ADMIN_AUTH:1:")
	assert.NoError(t, err)

	assert.Equal(t, []string{"JWT_ADMIN_AUTH:1:abc", "JWT_ADMIN_AUTH:1:def", "JWT_ADMIN_AUTH:1:ghi"}, keys)
}
