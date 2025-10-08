package subscriber

import (
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
)

type EtcdConf subscriber.EtcdConf

// MustNewEtcdSubscriber returns an etcd Subscriber, exits on errors.
func MustNewEtcdSubscriber(conf EtcdConf) subscriber.Subscriber {
	s, err := subscriber.NewEtcdSubscriber(subscriber.EtcdConf(conf))
	logx.Must(err)
	return s
}
