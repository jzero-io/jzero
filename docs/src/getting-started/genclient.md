---
title: 生成客户端代码
icon: clarity:thin-client-line
order: 5
---

## 生成 Swagger 文档

### 使用方法

::: code-tabs#shell

@tab jzero cli

```bash
cd your_project
jzero gen swagger
```

@tab jzero Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
```
:::

**Swagger UI 地址**: `localhost:8001/swagger`

## 生成 Zrpc 客户端

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen zrpcclient --name simplerpcclient
```

@tab Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen zrpcclient
```
:::

**代码示例**:

```go
package main

import (
	"context"
	"fmt"
	"simplerpc/simlerpcclient"
	"simplerpc/simlerpcclient/typed/version"

	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	cli, err := zrpc.NewClientWithTarget("localhost:8001")
	if err != nil {
		panic(err)
	}
	clientset, err := simlerpcclient.NewClientset(cli)
	if err != nil {
		panic(err)
	}
	versionResponse, err := clientset.Version().Version(context.Background(), &version.VersionRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(versionResponse)
}
```