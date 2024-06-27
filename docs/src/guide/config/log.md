---
title: 日志配置
icon: gears
star: true
order: 1
category: 配置
tag:
  - Guide
---

修改 etc/etc.yaml

```yaml
Log:
  KeepDays: 30
  Level: info
  MaxBackups: 7
  MaxSize: 50
  Mode: file
  Rotation: size
  ServiceName: app1
  encoding: plain
```

默认配置下日志最大占用空间: 2G

计算规则如下: 

logs 文件夹一共 5 个文件. 每个文件最大占用 50MB, 最多备份 7 个. 即 50MB * 8 * 5

## logtoconsole

在 go-zero 中, 设置日志 mode 为 file 或者 volume 时, 无法在控制台上查看日志, 解决办法: 使用 jzero-contrib 下的方法: logtoconsole.Must() 即可

```go
package main

import (
	"github.com/jzero-io/jzero-contrib/logtoconsole"
	"github.com/jzero-io/jzero-contrib/swaggerv2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	logConf := logx.LogConf{
		Mode:     "file",
		Path:     "logs",
		Encoding: "plain",
	}
	server := rest.MustNewServer(rest.RestConf{
		Port: 8001,
		ServiceConf: service.ServiceConf{
			Log: logConf,
		},
	})
	logtoconsole.Must(logConf)
	swaggerv2.RegisterRoutes(server, swaggerv2.WithSwaggerPath("docs"))

	logx.Info("starting server")
	server.Start()
}
```


