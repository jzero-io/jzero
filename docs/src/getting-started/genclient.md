---
title: 生成客户端代码
icon: clarity:thin-client-line
order: 5
---

## 生成 Swagger

### 使用方法

::: code-tabs#shell

@tab jzero cli

```bash
cd your_project
jzero gen swagger
# 合并成一个 swagger.json 文件
jzero gen swagger --merge
```

@tab jzero Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# 合并成一个 swagger.json 文件
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger --merge
```
:::

**Swagger 文件**: `desc/swagger`

**Swagger UI 地址**: `localhost:8001/swagger`

## 生成 Zrpc 客户端

生成 zrpc 客户端存在以下场景:

* 直接对服务端的 proto 文件生成客户端代码, 与主服务共用一个 go module, 其他服务要引用时, 需要引用整个源码才能使用(不推荐)
* 将服务端的 proto 文件复制到要引用的服务, 使用 `jzero gen zrpcclient` 生成到当前项目中(简单场景下推荐使用)
* 直接对服务端的 proto 文件生成客户端代码, 有独立的 go module, 通过 ci 流程推送到单独的远程仓库, 其他服务引用时直接 go get 引入(大型项目推荐)

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen zrpcclient --output simplerpcclient
# 设置 zrpcclient 为独立 go module
jzero gen zrpcclient --output simplerpcclient --goModule gitlab.xx.com/xx/simplerpcclient
```

@tab Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen zrpcclient --output simplerpcclient
# 设置 zrpcclient 为独立 go module
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen zrpcclient --output simplerpcclient --goModule gitlab.xx.com/xx/simplerpcclient
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

### 服务端调用其他 zrpcclient 场景

推荐将生成的 zrpcclient 放在 third_party 下, 并且为了区分服务端和客户端, 建议使用单独的 .jzero.yaml 生成 zrpcclient

如:

`third_party/matchingclient/.jzero.yaml`

```yaml
gen:
  zrpcclient:
    output: .
```

在 `third_party/matchingclient` 下执行 `jzero gen zrpcclient`

```shell
$ tree third_party/matchingclient -a
third_party/matchingclient
├── .jzero.yaml
├── clientset.go
├── desc
│   └── proto
│       └── matching.proto
├── model
│   └── types
│       └── matching
│           ├── matching.pb.go
│           └── matching_grpc.pb.go
└── typed
    └── matching
        └── matching.go
```

