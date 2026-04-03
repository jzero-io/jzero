---
title: "使用 jzero 一键搭建 AI 基座项目"
icon: /icons/emojione-v1-rocket.svg
description: "通过标准化框架 + AI 辅助业务逻辑的方式，快速搭建企业级 Go 项目"
---

## 为什么需要 AI 基座项目？

### 为什么需要 AI 基座项目？

- ✅ **统一框架**：避免重复选型和技术债务，确保长期可维护性
- ✅ **标准化结构**：清晰的分层架构、统一的代码组织和工程实践，解决代码混乱、维护困难的问题
- ✅ **质量与效率**：开箱即用的日志、监控、链路追踪、限流熔断等生产级基础设施，减少 80% 重复工作
- ✅ **团队规范**：内置 Go 最佳实践，AI 友好的架构设计，降低学习成本，提升协作效率

> 💡 **核心价值**：通过**标准化框架 + AI 辅助业务逻辑**，你不需要重复造轮子，AI 也不需要"猜测"框架规范，只需专注于业务逻辑本身！

## 快速开始

### 创建项目

jzero 提供了多种项目模板，满足不同场景需求：

#### 1️⃣ API 项目（RESTful 服务）

```bash
jzero new myproject --frame api
```

**适用场景**：
- ✅ RESTful API 服务
- ✅ 前后端分离的后端服务
- ✅ 移动端后端服务
- ✅ 微服务中的 API 网关

**特点**：
- 🚀 轻量级高性能
- 📝 自动生成 Swagger 文档
- 🔒 内置 JWT 认证
- 🎯 完善的请求验证

#### 2️⃣ RPC 项目（gRPC 服务）

```bash
jzero new myproject --frame rpc
```

**适用场景**：
- ✅ 微服务间通信
- ✅ 高性能内部服务
- ✅ 需要强类型的接口

**特点**：
- ⚡ 高性能二进制传输
- 🔄 自动负载均衡
- 📡 服务发现
- 🔐 内置重试和熔断

#### 3️⃣ Gateway 项目（API 网关）

```bash
jzero new myproject --frame gateway
```

**适用场景**：
- ✅ 统一入口管理
- ✅ 多服务聚合
- ✅ 需要同时支持 HTTP 和 gRPC

**特点**：
- 🌐 同时支持 HTTP 和 gRPC
- 🔄 智能路由转发
- ⚖️ 负载均衡
- 🛡️ 统一认证鉴权

### 启动服务

```bash
cd myproject

# 下载依赖
go mod tidy

# 启动服务
go run main.go server

# 访问 Swagger UI
open http://localhost:8001/swagger
```

## 按需启用：可选功能模块

你的项目可能需要数据库、缓存等功能，jzero 支持**按需启用**：

```bash
# 需要数据库 + Redis 缓存
jzero new myproject --features model,redis

# 需要数据库 + 数据库缓存
jzero new myproject --features model,cache

# 需要数据库 + Redis + 数据库缓存
jzero new myproject --features model,redis,cache
```

**功能说明**：

| 功能 | 说明 | 适用场景 |
|------|------|---------|
| `model` | 关系型数据库支持 | 需要持久化存储 |
| `redis` | Redis 缓存支持 | 需要高性能缓存 |
| `cache` | 数据库查询缓存 | 减轻数据库压力 |

## 🏢 自定义企业级模板

jzero 默认模板包含标准化的框架代码和最佳实践，但企业内部通常需要添加：
- 🔧 **CI/CD 流水线配置**（Jenkins）
- 📊 **企业监控告警**（Prometheus、Grafana、日志规范）
- 🛡️ **安全合规要求**（鉴权、审计、加密）
- 🏗️ **中间件集成**（MQ、OSS、第三方服务）

jzero 支持自定义企业模板，让所有项目都符合企业标准：

### 使用自定义模板

```bash
# 使用企业自定义模板创建项目
jzero new myproject --template https://github.com/your-org/jzero-template --branch base
```

## 从基座到业务：AI 辅助开发实战

有了 AI 基座项目后，你只需要专注于业务逻辑开发。让我们看看如何快速开发一个用户管理功能：

### 场景：开发用户注册功能

**使用 jzero + AI Skills**：

```bash
# 只需要一句话
"用 jzero-skills 创建用户注册功能，支持用户名、邮箱、密码，需要验证和去重"
```

**AI 自动完成**：
```
✅ 生成 desc/api/user.api
✅ 生成 desc/sql/user.sql
✅ 执行 jzero gen 生成框架代码
✅ 实现 Logic 层业务逻辑
✅ 包含完整的验证和错误处理
```

## 总结

jzero 通过**标准化框架 + AI 辅助**的方式，让 Go 项目开发变得前所未有的简单：

1. **快速启动**：一行命令创建包含最佳实践的项目基座
2. **灵活定制**：支持自定义企业模板，确保团队规范统一
3. **AI 增强**：AI 理解框架规范，专注业务逻辑实现

**觉得有用？请给 jzero 一个 ⭐ Star 支持我们持续改进！**

GitHub: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)
Jzero 官网: [https://jzero.io](https://jzero.io)