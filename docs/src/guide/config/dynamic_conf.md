---
title: 动态配置
icon: file-icons:config-go
order: 1.1
star: true
category: 配置
tag:
  - Guide
---

:::tip jzero >= v0.30.0+ 默认模板均支持动态配置特性
:::

* 基于 fsnotify 监听配置文件变化, 支持动态配置, 配置文件支持 yaml, json, toml.
* 使用 envsubst 库支持从环境变量读取值, 注意使用方式, 见下文.

## 使用 fsnotify 实现 go-zero 的 subscriber.Subscriber 接口

```go
package dynamic_conf

import (
	"os"
	"path/filepath"

	"github.com/a8m/envsubst"
	"github.com/eddieowens/opts"
	"github.com/fsnotify/fsnotify"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
)

var _ subscriber.Subscriber = (*FsNotify)(nil)

type FsNotify struct {
	path string

	// options
	options FsNotifyOpts

	*fsnotify.Watcher
}

type FsNotifyOpts struct {
	UseEnv bool
}

func (opts FsNotifyOpts) DefaultOptions() FsNotifyOpts {
	return FsNotifyOpts{}
}

func WithUseEnv(useEnv bool) opts.Opt[FsNotifyOpts] {
	return func(o *FsNotifyOpts) {
		o.UseEnv = useEnv
	}
}

func NewFsNotify(path string, op ...opts.Opt[FsNotifyOpts]) (*FsNotify, error) {
	o := opts.DefaultApply(op...)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FsNotify{
		path:    path,
		Watcher: watcher,
		options: o,
	}, nil
}

func (f *FsNotify) AddListener(listener func()) error {
	go func() {
		for {
			select {
			case event, ok := <-f.Watcher.Events:
				if !ok {
					return
				}
				if (event.Has(fsnotify.Write) || event.Has(fsnotify.Rename)) &&
					filepath.ToSlash(filepath.Clean(event.Name)) == filepath.Clean(filepath.ToSlash(f.path)) {
					logx.Infof("listen %s %s event", event.Name, event.Op)
					listener()
				}
			case err, ok := <-f.Watcher.Errors:
				if !ok {
					return
				}
				logx.Errorf("error: %v", err)
			}
		}
	}()

	// see: https://github.com/fsnotify/fsnotify/issues/363
	if err := f.Watcher.Add(filepath.Dir(f.path)); err != nil {
		return err
	}
	return nil
}

func (f *FsNotify) Value() (string, error) {
	file, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	if f.options.UseEnv {
		file, err = envsubst.Bytes(file)
		if err != nil {
			return "", err
		}
	}

	return string(file), nil
}
```

## 开启环境变量 dynamic_conf.WithUseEnv

```go
ss, err := dynamic_conf.NewFsNotify(cfgFile, dynamic_conf.WithUseEnv(true))
logx.Must(err)
cc := configurator.MustNewConfigCenter[config.Config](configurator.Config{
	Type: "yaml",
	}, ss)
c, err := cc.GetConfig()
logx.Must(err)
```

其中配置文件内容如下:

* 将需要从环境变量读取值的配置使用 ${} 包裹
* 原理为: 先读取配置文件的内容, 然后使用 envsubst 将 ${} 包裹的值替换为环境变量的值

```yaml
name: jzero
databaseType: ${DatabaseType}
```

## 替换 fsnotify 为 etcd 等注册中心

> https://go-zero.dev/docs/tasks/configcenter?_highlight=configcenter

```go
ss := subscriber.MustNewEtcdSubscriber(subscriber.EtcdConf{
        Hosts: []string{"localhost:2379"}, // etcd 地址
        Key:   "test1",    // 配置key
    })

// 创建 configurator
cc := configurator.MustNewConfigCenter[TestSt](configurator.Config{
        Type: "yaml",
    }, ss)
```
