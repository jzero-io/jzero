package redistest

import (
	"testing"

	"github.com/alicebob/miniredis/v2"

	"github.com/jzero-io/jzero/core/stores/redis"
)

// CreateRedis returns an in process redis.Redis.
func CreateRedis(t *testing.T) *redis.Redis {
	t.Helper()
	r, _ := CreateRedisWithClean(t)
	return r
}

// CreateRedisWithClean returns an in process redis.Redis and a clean function.
func CreateRedisWithClean(t *testing.T) (r *redis.Redis, clean func()) {
	t.Helper()
	mr := miniredis.RunT(t)
	return redis.New(mr.Addr()), mr.Close
}
