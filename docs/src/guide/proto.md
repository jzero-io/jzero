---
title: proto guide
icon: /icons/vscode-icons-file-type-protobuf.svg
star: true
order: 1
---

## Features

- ✅ **Support multiple proto files**: Can define multiple proto files in project (e.g., user.proto, order.proto, product.proto)
- ✅ Support **importing common proto** files
- ✅ **One-click generate RPC client**: Generate independent RPC client code, decouple from server dependency
- ✅ **Built-in field validation**: Automatic parameter validation based on `buf.validate`
- ✅ **Flexible middleware configuration**: Support configuring HTTP/RPC middleware for entire service or single method

## proto file example

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

## proto field validation

see: [protovalidate](https://buf.build/docs/protovalidate/)

```protobuf
syntax = "proto3";

package versionpb;

import "buf/validate/validate.proto";

option go_package = "./pb/versionpb";

// Use built-in rules
message CreateRequest {
  string email = 1 [(buf.validate.field).string.email = true];
  int32 age = 2 [(buf.validate.field).int32.gt = 17];
  string name = 3 [(buf.validate.field).string.min_len = 2];
}

// cel expression, supports custom message
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

**If built-in rules need internationalization, can use [protovalidate-translator](https://github.com/jzero-io/protovalidate-translator)**

## middleware

Add middleware, separate multiple middleware with commas

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

Detailed explanation:
* option (jzero.api.http_group) adds http middleware to all methods under this service
* option (jzero.api.http) only adds http middleware to specific method
* option (jzero.api.zrpc_group) adds zrpc middleware to all methods under this service
* option (jzero.api.zrpc) only adds zrpc middleware to specific method

After executing `jzero gen`, following files will be generated, using auth as example:
* internal/middleware/authmiddleware.go
* internal/middleware/middleware_gen.go
