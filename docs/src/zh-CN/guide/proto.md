---
title: proto 指南
icon: vscode-icons:file-type-protobuf
star: true
order: 1
---

## 特性

- ✅ **支持多 proto 文件**：可在项目中定义多个 proto 文件（如 user.proto、order.proto、product.proto）
- ✅ 支持**引入公共 proto** 文件
- ✅ **一键生成 RPC 客户端**：生成独立的 RPC 客户端代码，脱离服务端依赖，解耦服务端和客户端
- ✅ **内置字段验证**：基于 `buf.validate` 实现自动参数校验
- ✅ **灵活中间件配置**：支持为整个 service 或单个 method 配置 HTTP/RPC 中间件

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

// 使用内置规则
message CreateRequest {
  string email = 1 [(buf.validate.field).string.email = true];
  int32 age = 2 [(buf.validate.field).int32.gt = 17];
  string name = 3 [(buf.validate.field).string.min_len = 2];
}

// cel 表达式, 支持自定义 message
message GetRequest {
  int32 id = 1 [
    (buf.validate.field).cel = {
      id: "id.length"
      message: "id must be greater than 0 and less than 100000"
      expression: "this > 0 && this < 100000"
    }
  ];
}
```

**内置规则如果有国际化需求, 可使用 [protovalidate-translator](https://github.com/jzero-io/protovalidate-translator)**

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