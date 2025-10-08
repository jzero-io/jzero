---
title: proto 使用文档
icon: vscode-icons:file-type-protobuf
star: true
order: 1
category: 开发
tag:
  - Guide
---

:::tip jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).

jzero 在自动生成代码的时候会自动识别 desc/proto 下的文件并自动注册到 zrpc 上.
jzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate.
:::

jzero 框架的理念是:

* 不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 machine.proto.
* 每个 proto 文件可以有多个 service. 对于复杂模块可以使用多个 service.

jzero 中 proto 规范:

* proto 文件引用规范: 依据于 go-zero 的 proto 规范， 即 service 的 rpc 方法中入参和出参的 proto 不能是 import 的 proto 文件中的 message, 只能在当前文件

## proto 文件示例

::: code-tabs

@tab credential.proto

```protobuf
syntax = "proto3";

package credentialpb;

import "google/api/annotations.proto";

option go_package = "./pb/credentialpb";

message Empty {}

message CredentialVersionResponse {
  string version = 1;
}

message CreateCredentialRequest {
  string name = 1;
  string type = 2;
}

message CreateCredentialResponse {
  string name = 1;
  string type = 2;
}

service credential {
  rpc CredentialVersion(Empty) returns(CredentialVersionResponse) {
    option (google.api.http) = {
      get: "/api/v1.0/credential/version"
    };
  };

  rpc CreateCredential(CreateCredentialRequest) returns(CreateCredentialResponse) {
    option (google.api.http) = {
      post: "/api/v1.0/credential/create"
      body: "*"
    };
  }
}
```

@tab machine.proto

```protobuf
syntax = "proto3";

package machinepb;

import "google/api/annotations.proto";

option go_package = "./pb/machinepb";

message Empty {}

message MachineVersionResponse {
  string version = 1;
}

message CreateMachineRequest {
  string name = 1;
  string type = 2;
}

message CreateMachineResponse {
  string name = 1;
  string type = 2;
}

service credential {
  rpc MachineVersion(Empty) returns(MachineVersionResponse) {
    option (google.api.http) = {
      get: "/api/v1.0/machine/version"
    };
  };

  rpc CreateMachine(CreateMachineRequest) returns(CreateMachineResponse) {
    option (google.api.http) = {
      post: "/api/v1.0/machine/create"
      body: "*"
    };
  }
}
```

@tab chain.proto(最复杂场景 proto 多 service)

```protobuf
syntax = "proto3";

package chainpb;

import "google/api/annotations.proto";

option go_package = "./pb/chainpb";

message Empty {}

message CreateNodeRequest {
  string name = 1;
  string type = 2;
}

message CreateNodeResponse {
  string name = 1;
  string type = 2;
}

message CreateNamespaceRequest {
  string name = 1;
  string type = 2;
}

message CreateNamespaceResponse {
  string name = 1;
  string type = 2;
}

service node {
  rpc CreateNode(CreateNodeRequest) returns(CreateNodeResponse) {
    option (google.api.http) = {
      post: "/api/v1.0/chain/node/create"
      body: "*"
    };
  }
}

service namespace {
  rpc CreateNamespace(CreateNamespaceRequest) returns(CreateNamespaceResponse) {
    option (google.api.http) = {
      post: "/api/v1.0/chain/namespace/create"
      body: "*"
    };
  }
}
```
:::

## proto 字段校验

see: [protovalidate](https://buf.build/docs/protovalidate/)

```protobuf
syntax = "proto3";

package versionpb;

import "buf/validate/validate.proto";

option go_package = "./pb/versionpb";

message VersionRequest {
  int32 id = 1 [
    (buf.validate.field).cel = {
      id: "id.length"
      message: "id must be greater than 0 and less than 100000"
      expression: "this > 0 && this < 100000"
    }
  ];
}
```

## proto 扩展

### middleware 的分组管理

:::tip 确保存在 desc/proto/jzero/api 文件夹

如果不存在, 请下载到本地 https://github.com/jzero-io/desc/tree/main/proto/jzero/api
:::

添加 middleware, 多个 middleware 使用逗号隔开

```protobuf
import "jzero/api/http.proto";
import "jzero/api/zrpc.proto";

service User {
    option (jzero.api.http_group) = {
        middleware: "auth",
    };

    rpc CreateUser(CreateUserRequest) returns(CreateUserResponse) {
        option (google.api.http) = {
            post: "/api/v1/user/create",
            body: "*"
        };
        option (jzero.api.zrpc) = {
            middleware: "withValue1",
        };
    };

    rpc ListUser(ListUserRequest) returns(ListUserResponse) {
        option (google.api.http) = {
            get: "/api/v1/user/{username}/list",
        };
    };
}
```

详细解释:
* option (jzero.api.http_group) 即将该 service 下的所有 method 都新增 http 中间件
* option (jzero.api.http) 只针对某个 method 新增 http 中间件
* option (jzero.api.zrpc_group) 即将该 service 下的所有 method 都新增 zrpc 中间件
* option (jzero.api.zrpc) 只针对某个 method 新增 zrpc 中间件

执行 `jzero gen` 后将会生成一下文件, 以 auth 为例:
* internal/middleware/authmiddleware.go
* internal/middleware/middleware_gen.go

修改一下 cmd/server 代码, 将生成的文件注册到 server 中:

```go
zrpc := server.RegisterZrpc(svcCtx.Config, svcCtx)
gw := gateway.MustNewServer(svcCtx.Config.Gateway.GatewayConf, middleware.WithHeaderProcessor())
middleware.RegisterGen(zrpc, gw)
```
