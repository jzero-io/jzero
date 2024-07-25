---
title: 生成客户端 SDK
icon: carbon:sdk
order: 5
---

## 自动生成 go http sdk

* kubernetes client-go style
* 根据 api group 和 proto service 进行业务分组
* 链式调用, 调用逻辑清晰
* 可自定义接口进行扩展

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen sdk
cd your_project_go
go mod tidy
```

@tab Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app jaronnie/jzero:latest gen sdk

cd your_project-go
go mod tidy
```
:::

### 使用实例

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/jzero-restc"
	"your_project-go"
	"your_project-go/model/your_project/types"
)

func main() {
	clientset, err := your_project.NewClientWithOptions(
		restc.WithAddr("127.0.0.1"),
		restc.WithPort("8001"),
		restc.WithProtocol("http"))
	if err != nil {
		panic(err)
	}

	result, err := clientset.Hello().SayHello(context.Background(), &types.SayHelloRequest{
		Message: "jzero",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Message)
}
```

## 自动生成 ts http sdk

:::tip working...
:::

## 自动生成 zrpc client

```shell
jzero gen zrpcclient
```

### 使用实例一: 直连 rpc 服务

```go
package main

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/zrpc"
	your_project "your_project/zrpcclient-go"
	"your_project/zrpcclient-go/model/your_project/pb/hellopb"
)

func main() {
	target, err := zrpc.NewClientWithTarget("localhost:8000")
	if err != nil {
		panic(err)
	}
	var cs your_project.Interface

	cs = your_project.NewClientset(your_project.WithYour_projectClient(target))

	hello, err := cs.Your_project().Hello().SayHello(context.Background(), &hellopb.SayHelloRequest{
		Message: "hello",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(hello)
}
```

### 使用实例二: 基于 etcd 连接

使用之前请修改服务端配置 etc/etc.yaml, [请查看 etcd 配置](config/etcd.md), 然后重启程序

```go
package main

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/zrpc"
	your_project "your_project/zrpcclient-go"
	"your_project/zrpcclient-go/model/your_project/pb/hellopb"
)

func main() {
	target, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "your_project.rpc",
		},
	})
	var cs your_project.Interface

	cs = your_project.NewClientset(your_project.WithYour_projectClient(target))

	hello, err := cs.Your_project().Hello().SayHello(context.Background(), &hellopb.SayHelloRequest{
		Message: "hello",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(hello)
}
```

