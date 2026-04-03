---
title: Generate client code
icon: /icons/clarity-thin-client-line.svg
order: 5
---

## Generate Swagger

### Usage

::: code-tabs#shell

@tab jzero cli

```bash
cd your_project
jzero gen swagger
# merge into single swagger.json file
jzero gen swagger --merge
```

@tab jzero Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# merge into single swagger.json file
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger --merge
```
:::

**Swagger files**: `desc/swagger`

**Swagger UI address**: `localhost:8001/swagger`

## Generate Zrpc client

There are several scenarios for generating zrpc client:

* Directly generate client code from server's proto file, share same go module with main service, other services need to reference entire source code to use (not recommended)
* Copy server's proto file to referencing service, use `jzero gen zrpcclient` to generate to current project (recommended for simple scenarios)
* Directly generate client code from server's proto file, has independent go module, push to separate remote repository through CI process, other services directly go get to reference (recommended for large projects)

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen zrpcclient --output simplerpcclient
# set zrpcclient as independent go module
jzero gen zrpcclient --output simplerpcclient --goModule gitlab.xx.com/xx/simplerpcclient
```

@tab Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen zrpcclient --output simplerpcclient
# set zrpcclient as independent go module
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen zrpcclient --output simplerpcclient --goModule gitlab.xx.com/xx/simplerpcclient
```
:::

**Code example**:

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

### Server calling other zrpcclient scenarios

Recommend placing generated zrpcclient under third_party, and to distinguish server and client, recommend using separate .jzero.yaml to generate zrpcclient

For example:

`third_party/matchingclient/.jzero.yaml`

```yaml
gen:
  zrpcclient:
    output: .
```

Execute `jzero gen zrpcclient` under `third_party/matchingclient`

```shell
$ tree third_party/matchingclient -a
third_party/matchingclient
├── .jzero.yaml
├── clientset.go
├── desc
│   └── proto
│       └── matching.proto
├── model
│   └── types
│       └── matching
│           ├── matching.pb.go
│           └── matching_grpc.pb.go
└── typed
    └── matching
        └── matching.go
```

