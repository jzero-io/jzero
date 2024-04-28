---
title: 生成客户端 SDK
icon: lightbulb
order: 5
---

## 自动生成客户端 sdk

::: code-tabs#shell

@tab jzero

```bash
cd app1
jzero gensdk --module=github.com/jaronnie/app1-go --dir=app1-go
go mod tidy
```

@tab Docker(amd64)

```bash
cd app1
docker run --rm \
  -v ./app1:/app/app1 jaronnie/jzero:latest \
  gensdk --module=github.com/jaronnie/app1-go --dir=app1-go

go mod tidy
```

@tab Docker(arm64)

```bash
cd app1
docker run --rm \
  -v ./app1:/app/app1 jaronnie/jzero:latest-arm64 \
  gensdk --module=github.com/jaronnie/app1-go --dir=app1-go

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