---
title: proto 规范
icon: vscode-icons:file-type-protobuf
star: true
order: 1
category: 开发
tag:
  - Guide
---

:::tip jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).

jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上.
jzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate.
:::

jzero 框架的理念是:

* 不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 machine.proto.
* 每个 proto 文件可以有多个 service. 对于复杂模块可以使用多个 service.
* 如需对模块进行版本管理, 应该是 credential.proto, credential_v2.proto 规范.

proto 规范:

* 依据于 go-zero 的 proto 规范. 即 service 的 rpc 方法中 入参和出参的 proto 不能是 import 的 proto 文件

规范文件实例:

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
