---
title: "Free Your Hands! jzero Boosts Go Development Efficiency 10x"
icon: streamline-ultimate:blog-blogger-logo
---

As a developer, have you ever encountered these problems:

- Having to repeatedly set up basic infrastructure every time you create a new project?
- Business code mixed with infrastructure code, difficult to maintain?
- Different team members have different coding styles, making review costly?
- Want unified development standards but don't know where to start?
- As the project grows, module decoupling and collaboration become increasingly difficult?

If you have these concerns, today's article is a must-read!

---

## What is jzero?

**jzero** is an enhanced development tool based on the go-zero framework:

🏗️ **Generate basic framework code through templates**: Automatically generate framework code (api → api framework code, proto → proto framework code, sql/remote database address → model code) based on descriptor files

🤖 **Generate business code through Agent Skills**: Built-in jzero-skills enables AI to generate business logic code that follows best practices

**Core Value and Design Philosophy**:

- ✅ **Developer Experience First**: Provides a simple, easy-to-use, one-stop production-ready solution, one-click initialize api/rpc/gateway projects, minimal commands to generate basic framework code
- ✅ **AI Empowered**: Built-in jzero-skills enables AI to generate business logic code that follows best practices
- ✅ **Template-Driven**: Default generation follows best practices, supports custom templates, can build enterprise-specific foundation based on remote template repositories
- ✅ **Plugin Architecture**: Module layering, plugin design, smoother team collaboration
- ✅ **Built-in Components**: Includes common tools like cache, migrate, configcenter, condition
- ✅ **Ecosystem Compatible**: Doesn't modify go-zero, maintains ecosystem compatibility while addressing existing pain points and extending new features
- ✅ **Flexible Interface**: Doesn't depend on specific database/cache/config center, free choice based on actual needs

---

GitHub: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)

Docs: [https://docs.jzero.io](https://docs.jzero.io)

## Basic Framework Code Generation

Automatically generate basic framework code based on describable files:

### api → api framework code

```go
info (
    go_package: "user" // Define generated type folder location
)

type User {
    id int `json:"id"`
    username string `json:"username"`
}

type PageRequest {
    page int `form:"page"`
    size int `form:"size"`
}

type PageResponse {
    total uint64 `json:"total"`
    list  []User `json:"list"`
}

@server (
    prefix: /api/user        // Route prefix
    group: user              // Generated handler/logic folder location
    jwt: JwtAuth             // Enable JWT authentication
    middleware: AuthX        // Middleware
    compact_handler: true    // Merge this group's handlers into one file
)
service userservice {
    @doc "User pagination"
    @handler Page
    get /page (PageRequest) returns (PageResponse)
}
```

→ Generate Handler, Logic, Types, route registration, middleware, etc.

**Feature Description**:
- ✅ `go_package` - Define types folder location, avoid types.go being too large
- ✅ `compact_handler: true` - Merge handlers of the same group into one file, reduce file count

### proto → rpc framework code

```proto
syntax = "proto3";

package user;
option go_package = "./types/user";

// Import jzero extensions
import "jzero/api/http.proto";
import "jzero/api/zrpc.proto";

import "google/api/annotations.proto";

// Import common proto
import "common/common.proto";

// Import validation rules
import "buf/validate/validate.proto";

message GetUserRequest {
  int64 id = 1;
}

message CreateUserRequest {
  string username = 1 [
    (buf.validate.field).string = {
      min_len: 3,
      max_len: 20,
      pattern: "^[a-zA-Z0-9_]+$"
    }
  ];
  string email = 2 [
    (buf.validate.field).string.email = true,
    (buf.validate.field).string.max_len = 254,
    (buf.validate.field).string.min_len = 3
  ];
  string password = 3 [
    (buf.validate.field).cel = {
      id: "password.length"
      message: "password must contain at least 8 characters"
      expression: "this.size() >= 8"
    }
  ];
}

message CreateUserResponse {
  int64 id = 1;
  string username = 2;
}

message GetUserResponse {
  int64 id = 1;
  string username = 2;
}

service UserService {
  // Add HTTP middleware for entire service
  option (jzero.api.http_group) = {
    middleware: "auth,log",
  };

  // Add RPC middleware for entire service
  option (jzero.api.zrpc_group) = {
    middleware: "trace",
  };

  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
    option (google.api.http) = {
      post: "/api/v1/user/create",
      body: "*"
    };
  }

  rpc GetUser(GetUserRequest) returns (GetUserResponse) {
    option (google.api.http) = {
      get: "/api/v1/user/{id}",
    };
  };
}
```

→ Generate RPC server code, client code, HTTP Gateway, middleware

**Feature Description**:
- ✅ **Support multiple proto files**: Can define multiple proto files in project (e.g., user.proto, order.proto, product.proto)
- ✅ Support **importing common proto** files
- ✅ **One-click generate RPC client**: Generate independent RPC client code, decouple from server, separate server and client
- ✅ **Built-in field validation**: Automatic parameter validation based on `buf.validate`, supports CEL expressions
- ✅ **Flexible middleware configuration**: Support configuring HTTP/RPC middleware for entire service or single method


### sql/remote database → Model code

```sql
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

→ Generate Model layer code, CRUD operations, supports complex queries

**Feature Description**:
- ✅ **Multiple data sources**: Support generating model code based on sql files or remote database connections
- ✅ **Auto-generate CRUD interfaces**: Automatically generate basic operations like create, read, update, delete
- ✅ **Complex query support**: Provide powerful chain queries for complex business scenarios
- ✅ **One code adapts to multiple databases**: Generated code compatible with MySQL, PostgreSQL, Sqlite and other databases, no need to regenerate, easily switch underlying database storage

**Flexible generation strategy**, greatly improving code generation efficiency for large projects:

```bash
# Only generate code for files changed in git
jzero gen --git-change

# Generate for specific file
jzero gen --desc desc/api/user.api
```

**Flexible configuration**, goodbye to complex commands:

Support multiple configuration methods in combination:
- ✅ Configuration file (.jzero.yaml)
- ✅ Command-line parameters
- ✅ Environment variables

```bash
# Default configuration .jzero.yaml
jzero gen

# Specify configuration file
jzero gen --config .jzero.dev.yaml
```

One-click switch between local development, testing, and production environments!

**Hooks Configuration**: Support executing custom scripts before and after code generation

```yaml
# .jzero.yaml

# Global hooks
hooks:
  before:
    - echo "Execute before jzero command"
  after:
    - echo "Execute after jzero command"

# gen command configuration
gen:
  hooks:
    before:
      - echo "Execute before generating code"
      - go mod tidy
    after:
      - echo "Execute after generating code"
```
---

## Generate Business Code Through Agent Skills

Based on jzero-skills, let AI automatically generate business code that follows best practices:

```bash
# Output AI Skills configuration to Claude (default ~/.claude/skills)
jzero skills init

# Output to current project
jzero skills init --output .claude/skills

# In Claude, describe requirements in natural language, recommend starting with jzero-skills
```

**What can AI do for you**:

**REST API Development**:
- ✅ Automatically write standard-compliant `.api` files (set `go_package`, `group`, `compact_handler`)
- ✅ Automatically execute `jzero gen --desc desc/api/xxx.api` to generate framework code
- ✅ Automatically implement Logic layer business logic, following Handler → Logic → Model three-layer architecture

**Database Operations**:
- ✅ Automatically create SQL migration files (xx.up.sql & xx.down.sql)
- ✅ Automatically execute database migration (`jzero migrate up`)
- ✅ Automatically generate Model code (`jzero gen --desc desc/sql/xxx.sql`)

**RPC Service Development**:
- ✅ Automatically write `.proto` files to define service interfaces
- ✅ Automatically generate RPC server and client code
- ✅ Automatically implement server business logic, following Handler → Logic → Model three-layer architecture

---

<video width="720" height="450" controls>
  <source src="https://oss.jaronnie.com/jzero-skills.mp4" type="video/mp4">
</video>


## Plugin Architecture

Support **plugin development**, loading functional modules as independent plugins:

```bash
# Create helloworld api service
jzero new helloword --frame api

cd helloworld

# Add api plugin
jzero new plugin_name --frame api --serverless

# Add api plugin (mono type, use helloworld's go module)
jzero new plugin_name_mono --frame api --serverless --mono

# Build and load all plugins
jzero serverless build

# Unload all plugins
jzero serverless delete

# Unload specific plugin
jzero serverless delete --plugin plugin_name
```

**Perfectly supports**:

- 📦 Functional module decoupling, independent development and testing
- 👥 Team collaboration, different teams responsible for different plugins
- 🔄 Load on demand, flexible assembly of functions

---

## Quick Experience, Get Started in 5 Minutes

```bash
# 1. Install jzero
go install github.com/jzero-io/jzero/cmd/jzero@latest

# 2. One-click environment check
jzero check

# 3. Create project
# api project
jzero new helloworld --frame api
# rpc project
jzero new helloworld --frame rpc
# gateway project
jzero new helloworld --frame gateway

cd helloworld

# Download dependencies
go mod tidy

# Run service
go run main.go server

# Built-in Swagger UI
# http://localhost:8001/swagger
```

---

## Related Ecosystem

### jzero-intellij IDE Plugin

If you are a **GoLand / IntelliJ IDEA** user, **jzero-intellij plugin** will greatly enhance your development experience!

**Core Features**:
- ✅ One-click create descriptor files api/proto/sql
- ✅ API file intelligent highlighting
- ✅ File navigation, jump between api/proto and logic files
- ✅ Descriptor file line header execution button to generate code
- ✅ Configuration file .jzero.yaml execution button to generate code

<video width="720" height="450" controls>
  <source src="https://oss.jaronnie.com/jzero-intellij.mp4" type="video/mp4">
</video>

**Download**: https://github.com/jzero-io/jzero-intellij/releases

### jzero-admin Backend Management System

Backend management system based on jzero, built-in RBAC permission management, ready to use

**Core Features**:
- ✅ Complete permission system (user/menu/role)
- ✅ Multi-database support (MySQL/PostgreSQL/SQLite)
- ✅ Backend plugin support
- ✅ Internationalization support

![](https://oss.jaronnie.com/image-20251217134305041.png)

![](https://oss.jaronnie.com/image-20251217134332958.png)

![](https://oss.jaronnie.com/image-20251217134400658.png)

**Online Demo**:

- Aliyun Function Compute: [https://jzero-admin.jaronnie.com](https://jzero-admin.jaronnie.com)
- Vercel: [https://admin.jzero.io](https://admin.jzero.io)

**GitHub**: [https://github.com/jzero-io/jzero-admin](https://github.com/jzero-io/jzero-admin)

# In Conclusion

**jzero's mission is to make Go development simpler and more efficient. If interested, join us to explore new possibilities in Go development!** 🎉

**Find it useful? Please give jzero a ⭐ Star to support our continued improvement!**
