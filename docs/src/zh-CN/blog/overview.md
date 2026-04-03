---
title: 解放双手！jzero 让 Go 开发效率提升 10 倍
icon: /icons/streamline-ultimate-blog-blogger-logo.svg
---

作为一名开发者，你是否也遇到过这些问题：

- 每次新建项目都要重复搭建基础架构？
- 业务代码和基础设施代码混杂，难以维护？
- 团队成员代码风格各异，review 成本高？
- 想要统一的开发规范，却不知如何下手？
- 项目规模变大后，模块解耦和协作越来越困难？

如果你有以上困扰，那么今天的文章绝对不能错过！

---

## 什么是 jzero？

**jzero** 是基于 go-zero 框架开发的增强型开发工具：

🏗️ **通过模板生成基础框架代码**：基于描述文件自动生成框架代码（api → api 框架代码、proto → proto 框架代码、sql/远程数据库地址 → model 代码）

🤖 **通过 Agent Skills 生成业务代码**：内置 jzero-skills，让 AI 生成符合最佳实践的业务逻辑代码

**核心价值与设计理念**：

- ✅ **开发体验优先**：提供简单好用的一站式生产可用解决方案，一键初始化 api/rpc/gateway 项目，极简指令生成基础框架代码
- ✅ **AI 赋能**：内置 jzero-skills，让 AI 生成符合最佳实践的业务逻辑代码
- ✅ **模板驱动**：默认生成即最佳实践，支持自定义模板，可基于远程模板仓库打造企业专属底座
- ✅ **插件化架构**：模块分层、插件设计，团队协作更顺畅
- ✅ **内置开发组件**：包含缓存(cache)、数据库迁移(migrate)、配置中心(configcenter)、数据库查询(condition)等常用工具
- ✅ **生态兼容**：不修改 go-zero，保持生态兼容，解决已有痛点问题并扩展新功能
- ✅ **接口灵活**：不依赖特定数据库/缓存/配置中心，可根据实际需求自由选择

---

github 地址: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)

文档地址: [https://docs.jzero.io](https://docs.jzero.io)

## 基础框架代码生成

基于可描述文件自动生成基础框架代码：

### api → api 框架代码

```go
info (
    go_package: "user" // 定义生成的 type 文件夹位置
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
    prefix: /api/user        // 路由前缀
    group: user              // 生成的 handler/logic 文件夹位置
    jwt: JwtAuth             // 启用 JWT 认证
    middleware: AuthX        // 中间件
    compact_handler: true    // 合并该 group 的 handler 到同一个文件
)
service userservice {
    @doc "用户分页"
    @handler Page
    get /page (PageRequest) returns (PageResponse)
}
```

→ 生成 Handler、Logic、Types、路由注册、中间件等

**特性说明**：
- ✅ `go_package` - 定义 types 生成的文件夹位置，避免 types.go 过大
- ✅ `compact_handler: true` - 同一组的 handler 合并到同一个文件，减少文件数量

### proto → rpc 框架代码

```proto
syntax = "proto3";

package user;
option go_package = "./types/user";

// 引入 jzero 扩展
import "jzero/api/http.proto";
import "jzero/api/zrpc.proto";

import "google/api/annotations.proto";

// 引入公共 proto
import "common/common.proto";

// 引入验证规则
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
  // 为整个 service 添加 HTTP 中间件
  option (jzero.api.http_group) = {
    middleware: "auth,log",
  };

  // 为整个 service 添加 RPC 中间件
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

→ 生成 RPC 服务端代码、客户端代码、HTTP Gateway、中间件

**特性说明**：
- ✅ **支持多 proto 文件**：可在项目中定义多个 proto 文件（如 user.proto、order.proto、product.proto）
- ✅ 支持**引入公共 proto** 文件
- ✅ **一键生成 RPC 客户端**：生成独立的 RPC 客户端代码，脱离服务端依赖，解耦服务端和客户端
- ✅ **内置字段验证**：基于 `buf.validate` 实现自动参数校验，支持 CEL 表达式
- ✅ **灵活中间件配置**：支持为整个 service 或单个 method 配置 HTTP/RPC 中间件
 

### sql/远程数据库 → Model 代码

```sql
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

→ 生成 Model 层代码、CRUD 操作，支持复杂查询

**特性说明**：
- ✅ **多种数据源**：支持基于 sql 文件或远程数据库连接生成 model 代码
- ✅ **自动生成 CRUD 接口**：自动生成增删改查等基础操作
- ✅ **复杂查询支持**：提供强大的链式查询处理复杂业务场景
- ✅ **一套代码适配多数据库**：生成的代码兼容 MySQL、PostgreSQL、Sqlite 等多种数据库，无需重新生成，轻松切换数据库底层存储

**灵活生成策略**，极大提升大型项目代码生成效率：

```bash
# 只生成 git 改动的文件对应的代码
jzero gen --git-change

# 指定文件生成
jzero gen --desc desc/api/user.api
```

**灵活配置**，告别复杂指令：

支持多种配置方式自由组合：
- ✅ 配置文件（.jzero.yaml）
- ✅ 命令行参数
- ✅ 环境变量

```bash
# 默认配置 .jzero.yaml
jzero gen

# 指定配置文件
jzero gen --config .jzero.dev.yaml
```

本地开发、测试、生产环境一键切换！

**Hooks 配置**：支持在代码生成前后执行自定义脚本

```yaml
# .jzero.yaml

# 全局 hooks
hooks:
  before:
    - echo "执行 jzero 指令前执行"
  after:
    - echo "执行 jzero 指令后执行"

# gen 指令配置
gen:
  hooks:
    before:
      - echo "执行生成代码前执行"
      - go mod tidy
    after:
      - echo "执行生成代码后执行"
```
---

## 通过 Agent Skills 生成业务代码

基于 jzero-skills，让 AI 自动生成符合最佳实践的业务代码：

```bash
# 输出 AI Skills 配置到 Claude（默认 ~/.claude/skills）
jzero skills init

# 输出到当前项目
jzero skills init --output .claude/skills

# 在 Claude 中用自然语言描述需求, 推荐使用 jzero-skills 开头
```

**AI 能帮你做什么**：

**REST API 开发**：
- ✅ 自动编写符合规范的 `.api` 文件（设置 `go_package`、`group`、`compact_handler`）
- ✅ 自动执行 `jzero gen --desc desc/api/xxx.api` 生成框架代码
- ✅ 自动实现 Logic 层业务逻辑，遵循 Handler → Logic → Model 三层架构

**数据库操作**：
- ✅ 自动创建 SQL 迁移文件（xx.up.sql & xx.down.sql）
- ✅ 自动执行数据库迁移（`jzero migrate up`）
- ✅ 自动生成 Model 代码（`jzero gen --desc desc/sql/xxx.sql`）

**RPC 服务开发**：
- ✅ 自动编写 `.proto` 文件定义服务接口
- ✅ 自动生成 RPC 服务端和客户端代码
- ✅ 自动实现服务端业务逻辑，遵循 Handler → Logic → Model 三层架构

---

<video width="720" height="450" controls>
  <source src="https://oss.jaronnie.com/jzero-skills.mp4" type="video/mp4">
</video>


## 插件化架构

支持**插件化开发**，将功能模块作为独立插件加载：

```bash
# 创建 helloworld api 服务
jzero new helloword --frame api

cd helloworld

# 增加 api 插件
jzero new plugin_name --frame api --serverless

# 增加 api 插件(mono类型，即使用 helloworld 的 go module)
jzero new plugin_name_mono --frame api --serverless --mono

# 编译并加载所有插件
jzero serverless build

# 卸载所有插件
jzero serverless delete

# 卸载指定插件
jzero serverless delete --plugin plugin_name
```

**完美支持**：

- 📦 功能模块解耦，独立开发和测试
- 👥 团队协作，不同团队负责不同插件
- 🔄 按需加载，灵活组装功能

---

## 快速体验，5 分钟上手

```bash
# 1. 安装 jzero
go install github.com/jzero-io/jzero/cmd/jzero@latest

# 2. 一键检查环境
jzero check

# 3. 创建项目
# api 项目
jzero new helloworld --frame api
# rpc 项目
jzero new helloworld --frame rpc
# gateway 项目
jzero new helloworld --frame gateway

cd helloworld

# 下载依赖
go mod tidy

# 运行服务
go run main.go server

# 内置 Swagger UI
# http://localhost:8001/swagger
```

---

## 相关生态

### jzero-intellij IDE 插件

如果你是 **GoLand / IntelliJ IDEA** 用户，**jzero-intellij 插件**将极大提升你的开发体验！

**核心功能**：
- ✅ 一键创建可描述文件 api/proto/sql
- ✅ api 文件智能高亮
- ✅ 文件跳转，api/proto 与 logic 文件互相跳转
- ✅ 可描述文件行首执行按钮生成代码
- ✅ 配置文件 .jzero.yaml 增加执行按钮生成代码

<video width="720" height="450" controls>
  <source src="https://oss.jaronnie.com/jzero-intellij.mp4" type="video/mp4">
</video>

**下载地址**：https://github.com/jzero-io/jzero-intellij/releases

### jzero-admin 后台管理系统

基于 jzero 的后台管理系统，内置 RBAC 权限管理，开箱即用

**核心特性**：
- ✅ 完整权限系统(用户/菜单/角色)
- ✅ 多数据库支持(MySQL/PostgreSQL/SQLite)
- ✅ 后端插件化
- ✅ 国际化支持

![](https://oss.jaronnie.com/image-20251217134305041.png)

![](https://oss.jaronnie.com/image-20251217134332958.png)

![](https://oss.jaronnie.com/image-20251217134400658.png)

**在线演示**：

- 阿里云云函数：[https://jzero-admin.jaronnie.com](https://jzero-admin.jaronnie.com)
- Vercel：[https://admin.jzero.io](https://admin.jzero.io)

**GitHub**：[https://github.com/jzero-io/jzero-admin](https://github.com/jzero-io/jzero-admin)

# 写在最后

**jzero 的使命是让 Go 开发更简单、更高效。如果有兴趣，可以加入我们，一起探索 Go 开发的新可能！** 🎉

**觉得有用？也请给 jzero 一个 ⭐ Star，支持我们继续改进！**
