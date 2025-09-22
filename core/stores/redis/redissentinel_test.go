package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSentinelOptions(t *testing.T) {
	r := &Redis{}

	// Test WithMasterName
	WithMasterName("mymaster")(r)
	assert.Equal(t, "mymaster", r.masterName)
}

func TestNewRedisWithSentinel(t *testing.T) {
	conf := RedisConf{
		Host:       "sentinel1:26379,sentinel2:26379",
		Type:       SentinelType,
		MasterName: "mymaster",
		User:       "user",
		Pass:       "pass",
		Tls:        true,
		NonBlock:   true, // avoid connection test
	}

	redis, err := NewRedis(conf)
	assert.NoError(t, err)
	assert.NotNil(t, redis)
	assert.Equal(t, SentinelType, redis.Type)
	assert.Equal(t, "sentinel1:26379,sentinel2:26379", redis.Addr)
	assert.Equal(t, "mymaster", redis.masterName)
	assert.Equal(t, "user", redis.User)
	assert.Equal(t, "pass", redis.Pass)
	assert.True(t, redis.tls)
}

func TestSentinelValidation(t *testing.T) {
	tests := []struct {
		name      string
		conf      RedisConf
		expectErr bool
		errType   error
	}{
		{
			name: "valid sentinel config",
			conf: RedisConf{
				Type:       SentinelType,
				Host:       "sentinel1:26379,sentinel2:26379",
				MasterName: "mymaster",
			},
			expectErr: false,
		},
		{
			name: "invalid sentinel config - missing master name",
			conf: RedisConf{
				Type: SentinelType,
				Host: "sentinel1:26379",
			},
			expectErr: true,
			errType:   ErrEmptyMasterName,
		},
		{
			name: "invalid sentinel config - missing host",
			conf: RedisConf{
				Type:       SentinelType,
				MasterName: "mymaster",
			},
			expectErr: true,
			errType:   ErrEmptyHost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.conf.Validate()
			if test.expectErr {
				assert.Error(t, err)
				if test.errType != nil {
					assert.Equal(t, test.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
