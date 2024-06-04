---
title: 生成客户端 SDK
icon: code
order: 5
---

## 自动生成客户端 sdk

* kubernetes client-go style
* 根据 api group 和 proto service 进行业务分组
* 链式调用, 调用逻辑清晰
* 自带 fake client 支持单元测试
* 可自定义接口进行扩展

::: code-tabs#shell

@tab jzero

```bash
cd app1
jzero gen sdk --module=github.com/jaronnie/quickstart-go --dir=quickstart-go --scope quickstart
cd quickstart-go
go mod tidy
```

@tab Docker
```bash
cd app1
docker run --rm -v ${PWD}:/app/quickstart jaronnie/jzero:latest gen sdk --module=github.com/jaronnie/quickstart-go --dir=quickstart-go --scope quickstart -w quickstart

cd app1-go
go mod tidy
```
:::

## sdk 使用实例

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/jaronnie/app1-go"
	"github.com/jaronnie/app1-go/model/app1/types"
	"github.com/jaronnie/app1-go/rest"
)

func main() {
	clientset, err := app1.NewClientWithOptions(
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