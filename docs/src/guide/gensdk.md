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

## sdk 使用实例

```go
package main

import (
	"context"
	"fmt"
	
	"your_project-go"
	"your_project-go/model/quickstart/types"
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