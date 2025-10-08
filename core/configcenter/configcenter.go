package configcenter

import (
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
)

type Config configurator.Config

type ConfigCenter[T any] interface {
	configurator.Configurator[T]
	MustGetConfig() T
}

type configCenter[T any] struct {
	configurator.Configurator[T]
}

func (c *configCenter[T]) MustGetConfig() T {
	config, err := c.Configurator.GetConfig()
	logx.Must(err)
	return config
}

func (c *configCenter[T]) GetConfig() (T, error) {
	return c.Configurator.GetConfig()
}

func (c *configCenter[T]) AddListener(listener func()) {
	c.Configurator.AddListener(listener)
}

func MustNewConfigCenter[T any](c Config, subscriber subscriber.Subscriber) ConfigCenter[T] {
	cc, err := NewConfigCenter[T](c, subscriber)
	logx.Must(err)
	return cc
}

func NewConfigCenter[T any](c Config, subscriber subscriber.Subscriber) (ConfigCenter[T], error) {
	cc := &configCenter[T]{}

	center, err := configurator.NewConfigCenter[T](configurator.Config(c), subscriber)
	if err != nil {
		return nil, err
	}
	cc.Configurator = center
	return cc, nil
}
