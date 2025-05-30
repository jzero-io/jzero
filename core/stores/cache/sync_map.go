package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/errorx"
)

type (
	syncMapItem struct {
		data     []byte
		duration int64
	}

	syncMap struct {
		storage     *sync.Map
		errNotFound error
	}
)

func (sm *syncMap) ExpireCtx(ctx context.Context, key string, expire time.Duration) error {
	item, err := sm.read(key)
	if err != nil {
		return err
	}
	item.duration = time.Now().Unix() + int64(expire.Seconds())
	sm.storage.Store(key, item)
	return nil
}

func (sm *syncMap) GetPrefixKeysCtx(ctx context.Context, prefix string) ([]string, error) {
	var allKeys []string

	sm.storage.Range(func(key, value any) bool {
		keyStr, ok := key.(string)
		if !ok {
			return true
		}

		if len(keyStr) >= len(prefix) && keyStr[:len(prefix)] == prefix {
			allKeys = append(allKeys, keyStr)
		}
		return true
	})

	return allKeys, nil
}

func (sm *syncMap) SetNoExpireCtx(ctx context.Context, key string, val any) error {
	return sm.SetCtx(ctx, key, val)
}

// NewSyncMap creates an instance of SyncMap cache driver
func NewSyncMap(errNotFound error) Cache {
	return &syncMap{
		storage:     &sync.Map{},
		errNotFound: errNotFound,
	}
}

func (sm *syncMap) Del(keys ...string) error {
	return sm.DelCtx(context.Background(), keys...)
}

func (sm *syncMap) DelCtx(ctx context.Context, keys ...string) error {
	var be errorx.BatchError

	for _, key := range keys {
		if _, ok := sm.storage.Load(key); !ok {
			be.Add(sm.errNotFound)
		} else {
			sm.storage.Delete(key)
		}
	}

	return be.Err()
}

func (sm *syncMap) Get(key string, val any) error {
	return sm.GetCtx(context.Background(), key, val)
}

func (sm *syncMap) GetCtx(ctx context.Context, key string, val any) error {
	item, err := sm.read(key)
	if err == nil {
		return json.Unmarshal(item.data, val)
	}

	return sm.errNotFound
}

func (sm *syncMap) IsNotFound(err error) bool {
	return errors.Is(err, sm.errNotFound)
}

func (sm *syncMap) Set(key string, val any) error {
	return sm.SetCtx(context.Background(), key, val)
}

func (sm *syncMap) SetCtx(ctx context.Context, key string, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	sm.storage.Store(key, &syncMapItem{
		data:     data,
		duration: 0,
	})
	return nil
}

func (sm *syncMap) SetWithExpire(key string, val any, expire time.Duration) error {
	return sm.SetWithExpireCtx(context.Background(), key, val, expire)
}

func (sm *syncMap) SetWithExpireCtx(ctx context.Context, key string, val any, expire time.Duration) error {
	duration := int64(0)
	if expire > 0 {
		duration = time.Now().Unix() + int64(expire.Seconds())
	}
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	sm.storage.Store(key, &syncMapItem{data: data, duration: duration})
	return nil
}

func (sm *syncMap) Take(val any, key string, query func(val any) error) error {
	return sm.TakeCtx(context.Background(), val, key, query)
}

func (sm *syncMap) TakeCtx(ctx context.Context, val any, key string, query func(val any) error) error {
	if _, ok := sm.storage.Load(key); ok {
		return sm.GetCtx(ctx, key, val)
	}

	if err := query(val); err != nil {
		return err
	}

	return sm.SetCtx(ctx, key, val)
}

func (sm *syncMap) TakeWithExpire(val any, key string, query func(val any, expire time.Duration) error) error {
	return sm.TakeWithExpireCtx(context.Background(), val, key, query)
}

func (sm *syncMap) TakeWithExpireCtx(ctx context.Context, val any, key string, query func(val any, expire time.Duration) error) error {
	if _, ok := sm.storage.Load(key); ok {
		return sm.GetCtx(ctx, key, val)
	}

	// patch
	var expire time.Duration
	if value, ok := ctx.Value("expire").(time.Duration); ok {
		expire = value
	}

	if err := query(val, expire); err != nil {
		return err
	}

	return sm.SetWithExpireCtx(ctx, key, val, expire)
}

func (sm *syncMap) read(key string) (*syncMapItem, error) {
	v, ok := sm.storage.Load(key)
	if !ok {
		return nil, sm.errNotFound
	}

	item := v.(*syncMapItem)

	if item.duration == 0 {
		return item, nil
	}

	if item.duration <= time.Now().Unix() {
		sm.storage.Delete(key)
		return nil, errors.New("cache expired")
	}

	return item, nil
}
