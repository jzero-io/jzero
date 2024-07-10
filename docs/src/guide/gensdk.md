---
title: 生成客户端 SDK
icon: carbon:sdk
order: 5
---

## 自动生成 go http sdk

* kubernetes client-go style
* 根据 api group 和 proto service 进行业务分组
* 链式调用, 调用逻辑清晰
* 自带 fake client 支持单元测试
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
	
	"your_project-go"
	"your_project-go/model/your_project/types"
	"your_project-go/rest"
)

func main() {
	clientset, err := your_project.NewClientWithOptions(
		rest.WithAddr("127.0.0.1"),
		rest.WithPort("8001"),
		rest.WithProtocol("http"))
	if err != nil {
		panic(err)
	}

	result, err := clientset.Hello().HelloPathHandler(context.Background(), &types.PathRequest{
		Name: "jzero",
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
	"your_project/zrpcclient-go/hello"
	"your_project/zrpcclient-go/pb/hellopb"
)

func main() {
	target, err := zrpc.NewClientWithTarget("localhost:8000")
	if err != nil {
		panic(err)
	}

	logic := hello.NewHello(target)
	sayHello, err := logic.SayHello(context.Background(), &hellopb.SayHelloRequest{
		Message: "12345",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(sayHello.Message)
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
	"your_project/zrpcclient-go/hello"
	"your_project/zrpcclient-go/pb/hellopb"
)

func main() {
	client, err := zrpc.NewClient(zrpc.NewEtcdClientConf([]string{"127.0.0.1:2379"}, "your_project.rpc", "", ""))
	if err != nil {
		panic(err)
	}
	logic := hello.NewHello(client)

	sayHello, err := logic.SayHello(context.Background(), &hellopb.SayHelloRequest{
		Message: "12345",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(sayHello.Message)
}
```

