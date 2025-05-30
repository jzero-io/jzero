package subscriber

import (
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
)

// MustNewEtcdSubscriber returns an etcd Subscriber, exits on errors.
func MustNewEtcdSubscriber(conf subscriber.EtcdConf) subscriber.Subscriber {
	s, err := subscriber.NewEtcdSubscriber(conf)
	logx.Must(err)
	return s
}
