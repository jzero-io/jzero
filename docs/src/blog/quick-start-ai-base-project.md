---
title: "Build AI Base Project with jzero in One Click"
icon: emojione-v1:rocket
description: "Quickly build enterprise-grade Go projects through standardized framework + AI-assisted business logic"
---

## Why Do You Need an AI Base Project?

### Why Do You Need an AI Base Project?

- ✅ **Unified Framework**: Avoid repetitive technology selection and technical debt, ensure long-term maintainability
- ✅ **Standardized Structure**: Clear layered architecture, unified code organization and engineering practices, solving code chaos and maintenance difficulties
- ✅ **Quality and Efficiency**: Production-grade infrastructure including logging, monitoring, distributed tracing, rate limiting and circuit breaking out of the box, reducing 80% of repetitive work
- ✅ **Team Standards**: Built-in Go best practices, AI-friendly architecture design, reduced learning costs, improved collaboration efficiency

> 💡 **Core Value**: Through **standardized framework + AI-assisted business logic**, you don't need to reinvent the wheel, and AI doesn't need to "guess" framework specifications. Just focus on business logic itself!

## Quick Start

### Create Project

jzero provides various project templates to meet different scenario requirements:

#### 1️⃣ API Project (RESTful Service)

```bash
jzero new myproject --frame api
```

**Use Cases**:
- ✅ RESTful API services
- ✅ Backend services for frontend-backend separation
- ✅ Mobile backend services
- ✅ API gateways in microservices

**Features**:
- 🚀 Lightweight and high performance
- 📝 Auto-generated Swagger documentation
- 🔒 Built-in JWT authentication
- 🎯 Comprehensive request validation

#### 2️⃣ RPC Project (gRPC Service)

```bash
jzero new myproject --frame rpc
```

**Use Cases**:
- ✅ Inter-service communication in microservices
- ✅ High-performance internal services
- ✅ Services requiring strongly-typed interfaces

**Features**:
- ⚡ High-performance binary transmission
- 🔄 Automatic load balancing
- 📡 Service discovery
- 🔐 Built-in retry and circuit breaking

#### 3️⃣ Gateway Project (API Gateway)

```bash
jzero new myproject --frame gateway
```

**Use Cases**:
- ✅ Unified entry point management
- ✅ Multi-service aggregation
- ✅ Need to support both HTTP and gRPC

**Features**:
- 🌐 Support both HTTP and gRPC
- 🔄 Smart routing and forwarding
- ⚖️ Load balancing
- 🛡️ Unified authentication and authorization

### Start Service

```bash
cd myproject

# Download dependencies
go mod tidy

# Start service
go run main.go server

# Access Swagger UI
open http://localhost:8001/swagger
```

## On-Demand Enablement: Optional Feature Modules

Your project might need databases, caching, and other features. jzero supports **on-demand enablement**:

```bash
# Need database + Redis cache
jzero new myproject --features model,redis

# Need database + database cache
jzero new myproject --features model,cache

# Need database + Redis + database cache
jzero new myproject --features model,redis,cache
```

**Feature Description**:

| Feature | Description | Use Cases |
|---------|-------------|-----------|
| `model` | Relational database support | Need persistent storage |
| `redis` | Redis cache support | Need high-performance caching |
| `cache` | Database query cache | Reduce database pressure |

## 🏢 Custom Enterprise Templates

jzero default templates only contain basic framework code, but enterprises typically need to include content like CI/CD pipeline integration. jzero supports custom templates to add enterprise-specific content.

### Using Custom Templates

```bash
# Create project using enterprise custom template
jzero new myproject --template https://github.com/your-org/jzero-template --branch base
```

## From Base to Business: AI-Assisted Development in Action

With the AI base project, you only need to focus on business logic development. Let's see how to quickly develop a user management feature:

### Scenario: Develop User Registration Feature

**Using jzero + AI Skills**:

```bash
# Just one sentence
"Use jzero-skills to create user registration feature, supporting username, email, password with validation and deduplication"
```

**AI Automatically Completes**:

```
✅ Generate desc/api/user.api
✅ Generate desc/sql/user.sql
✅ Execute jzero gen to generate framework code
✅ Implement Logic layer business logic
✅ Include complete validation and error handling
```

## Summary

jzero makes Go project development easier than ever through **standardized framework + AI assistance**:

1. **Quick Start**: One command to create a project base with best practices
2. **Flexible Customization**: Support custom enterprise templates to ensure unified team standards
3. **AI Enhanced**: AI understands framework specifications and focuses on business logic implementation

**Find it useful? Please give jzero a ⭐ Star to support our continued improvement!**

GitHub: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)
Jzero Website: [https://jzero.io](https://jzero.io)