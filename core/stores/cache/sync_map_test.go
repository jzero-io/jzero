package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSyncMap(t *testing.T) {
	t.Run("TestSyncMap", func(t *testing.T) {
		cache := NewSyncMap(errors.New("not found"))

		err := cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:abc", "abc")
		assert.NoError(t, err)
		err = cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:def", "def")
		assert.NoError(t, err)
		err = cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:ghi", "ghi")
		assert.NoError(t, err)

		keys, err := cache.GetPrefixKeysCtx(context.Background(), "JWT_ADMIN_AUTH:1:")
		assert.NoError(t, err)
		assert.Contains(t, keys, "JWT_ADMIN_AUTH:1:abc")
		assert.Contains(t, keys, "JWT_ADMIN_AUTH:1:def")
		assert.Contains(t, keys, "JWT_ADMIN_AUTH:1:ghi")
	})
}

func TestSyncMapExpireCtx(t *testing.T) {
	r, err := miniredis.Run()
	assert.NoError(t, err)
	defer r.Close()

	cache := NewSyncMap(errors.New("not found"))

	err = cache.SetWithExpireCtx(context.Background(), "JWT_ADMIN_AUTH:1:abc", "abc", time.Duration(5)*time.Second)
	assert.NoError(t, err)

	var val any
	err = cache.Get("JWT_ADMIN_AUTH:1:abc", &val)
	assert.NoError(t, err)
	assert.Equal(t, "abc", val)

	time.Sleep(time.Second * 6)
	var newVal any
	err = cache.Get("JWT_ADMIN_AUTH:1:abc", &newVal)
	assert.Error(t, err)
}
