---
title: 生成客户端 SDK
icon: lightbulb
order: 5
---

## 自动生成客户端 sdk

* kubernetes client-go style
* 根据 api group 和 proto service 进行业务分组
* 链式调用, 调用逻辑清晰
* 自带 fake client 支持单元测试
* 可自定义接口进行扩展
* 支持基于统一 api 网关的多服务 sdk 集成在一个 sdk 库中

::: code-tabs#shell

@tab jzero

```bash
cd app1
jzero gensdk --module=github.com/jaronnie/app1-go --dir=app1-go
cd app1-go
go mod tidy
```

@tab Docker(amd64)

```bash
cd app1
docker run --rm \
  -v $PWD:/app/app1 jaronnie/jzero:latest \
  gensdk --module=github.com/jaronnie/app1-go --dir=app1-go -w app1

cd app1-go
go mod tidy
```

@tab Docker(arm64)

```bash
cd app1
docker run --rm \
  -v $PWD:/app/app1 jaronnie/jzero:latest-arm64 \
  gensdk --module=github.com/jaronnie/app1-go --dir=app1-go -w app1

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