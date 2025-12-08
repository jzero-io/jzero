---
title: proto 指南
icon: vscode-icons:file-type-protobuf
star: true
order: 1
---

:::tip jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).

jzero 在自动生成代码的时候会自动识别 desc/proto 下的文件并自动注册到 zrpc 上.
jzero 默认支持对 proto 的字段校验
:::

jzero 框架的理念是:

* 不同模块分在不同的 proto 文件下

jzero 中 proto 规范:

* proto 文件引用规范: 依据于 go-zero 的 proto 规范， 即 service 的 rpc 方法中入参和出参的 proto 不能是 import 的 proto 文件中的 message, 只能在当前文件

## proto 文件示例

```protobuf
syntax = "proto3";

package version;

import "google/api/annotations.proto";
import "grpc-gateway/protoc-gen-openapiv2/options/annotations.proto";

option go_package = "./types/version";

message VersionRequest {}

message VersionResponse {
  string version = 1;
  string goVersion = 2;
  string commit = 3;
  string date = 4;
}

service Version {
  rpc Version(VersionRequest) returns(VersionResponse) {
    option (google.api.http) = {
      get: "/version"
    };
  };
}
```

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

## middleware

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