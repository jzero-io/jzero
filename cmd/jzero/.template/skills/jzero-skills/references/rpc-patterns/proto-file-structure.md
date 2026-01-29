# Proto File Structure

## Overview

jzero supports multi-proto management (goctl native tool does not support this). When automatically generating code, jzero automatically recognizes files under `desc/proto/` and automatically registers them to zrpc.

jzero also supports proto field validation by default.

## jzero Framework Philosophy

**Different modules should be separated into different proto files**

## Proto File Standards

### Import Rule
Based on go-zero's proto standard: In service RPC methods, input and output parameters' proto messages cannot be from imported proto files - they must be defined in the current file only.

## Proto File Example

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

## File Structure

**Required elements:**

1. **Syntax declaration**: `syntax = "proto3";`
2. **Package name**: Unique identifier for the proto file
3. **Imports**:
   - `google/api/annotations.proto` - For HTTP mapping
   - `grpc-gateway/protoc-gen-openapiv2/options/annotations.proto` - For OpenAPI documentation
   - Custom imports as needed
4. **Go package option**: `option go_package = "./types/version";`
5. **Messages**: Request and response structures
6. **Service**: RPC method definitions with HTTP annotations

## Best Practices

### ✅ Correct Patterns

```protobuf
// Define request/response in the same file as the service
syntax = "proto3";

package user;

option go_package = "./types/user";

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

message CreateUserResponse {
  int32 id = 1;
}

service User {
  rpc CreateUser(CreateUserRequest) returns(CreateUserResponse);
}
```

### ❌ Incorrect Patterns

```protobuf
// DON'T: Define messages in separate imported files
syntax = "proto3";

package user;

import "common.proto";  // ❌ Messages should not be imported

option go_package = "./types/user";

service User {
  rpc CreateUser(common.Request) returns(common.Response);  // ❌ Wrong!
}
```

## Directory Structure

```
myproject/
├── desc/
│   └── proto/
│       ├── user.proto          # User service
│       ├── order.proto         # Order service
│       └── product.proto       # Product service
└── internal/
    ├── proto/                  # Generated proto code
    │   ├── user/
    │   ├── order/
    │   └── product/
    └── svc/
        └── servicecontext.go   # Auto-registers all proto services
```

## Code Generation

Generate RPC code from proto files:

```bash
# Generate from specific proto file
jzero gen --desc desc/proto/user.proto

# Generate all proto files
jzero gen
```

## Key Features

### Multi-Proto Support

Unlike goctl, jzero supports multiple proto files and automatically:
- Scans `desc/proto/` directory
- Generates code for all proto files
- Registers all services to zrpc server

### HTTP Gateway Support

Proto files can define HTTP mappings for REST endpoints:

```protobuf
service User {
  rpc CreateUser(CreateUserRequest) returns(CreateUserResponse) {
    option (google.api.http) = {
      post: "/api/v1/user/create"
      body: "*"
    };
  };

  rpc GetUser(GetUserRequest) returns(GetUserResponse) {
    option (google.api.http) = {
      get: "/api/v1/user/{id}"
    };
  };
}
```

### OpenAPI Documentation

Generate OpenAPI/Swagger documentation:

```protobuf
import "grpc-gateway/protoc-gen-openapiv2/options/annotations.proto";

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "User Service"
    version: "1.0"
    description: "User management service"
  }
};
```

## Related Topics

- [Proto Field Validation](proto-validation.md) - Adding validation to proto messages
- [Proto Middleware](proto-middleware.md) - Using middleware with proto services
- [jzero Documentation](https://docs.jzero.io) - Official documentation
