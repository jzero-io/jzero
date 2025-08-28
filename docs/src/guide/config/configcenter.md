---
title: 配置中心
icon: file-icons:config-go
order: 1.1
star: true
category: 配置
tag:
  - Guide
---

## 特性

* 支持基于 fsnotify 和 etcd 注册中心实现
* 默认使用 fsnotify 实现
* 支持自定义注册中心实现

## 默认使用 fsnotify 实现

```go
package main 

import (
	"fmt"
	
    "github.com/jzero-io/jzero/core/configcenter/subscriber"
    configurator "github.com/zeromicro/go-zero/core/configcenter"
)

type Config struct {
    Name string `json:"name"`
    DatabaseType string `json:"databaseType"`
}

func main() {
	cfgFile := "config.yaml"
	cc := configurator.MustNewConfigCenter[Config](configurator.Config{
		Type: "yaml",
	}, subscriber.MustNewFsnotifySubscriber(cfgFile, subscriber.WithUseEnv(true)))
	
	config , err := cc.GetConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
```

```yaml
name: jzero
databaseType: mysql
```

## 替换 fsnotify 为 etcd 等注册中心

> https://go-zero.dev/docs/tasks/configcenter?_highlight=configcenter

```go
package main

import (
	"fmt"

	subscriber "github.com/zeromicro/go-zero/core/configcenter/subscriber"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
)

type Config struct {
	Name string `json:"name"`
	DatabaseType string `json:"databaseType"`
}

func main() {
	cc := configurator.MustNewConfigCenter[Config](configurator.Config{
		Type: "yaml",
	},subscriber.MustNewEtcdSubscriber(subscriber.EtcdConf{
		Hosts: []string{"localhost:2379"}, // etcd 地址
		Key:   "test1",    // 配置key
	}))

	config , err := cc.GetConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
```
