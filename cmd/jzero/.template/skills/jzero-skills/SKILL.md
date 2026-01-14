---
name: jzero-skills
description: Comprehensive knowledge base for jzero framework (enhanced go-zero). Use this skill when working with jzero to understand correct patterns for REST APIs (Handler/Logic/Context architecture), RPC services (service discovery, load balancing), Gateway services, database operations (sqlx, MongoDB, caching), resilience patterns (circuit breaker, rate limiting), and jzero-specific features (git-change-based generation, flexible configuration, custom templates). Essential for generating production-ready jzero code that follows framework conventions.
license: Apache-2.0
---

# jzero Skills for AI Agents

Structured knowledge base optimized for AI agents to help developers work effectively with the [jzero](https://github.com/jzero-io/jzero) framework (enhanced go-zero).

## Overview

This skill provides AI agents with comprehensive jzero knowledge to:
- Generate accurate code following jzero conventions
- Understand the three-layer architecture (Handler → Logic → Model)
- Apply best practices for microservices development
- Use jzero-specific features (git-aware generation, flexible config)
- Build production-ready applications

## Quick Start

When helping with jzero development:

1. **For new projects**: Start with [Code Generation Workflow](#code-generation-workflow)
2. **For REST APIs**: Check [REST API File Structure](references/rest-api-patterns/api-file-structure.md) - ⚠️ Critical rules
3. **For databases**: Review [Database Best Practices](references/database-patterns/best-practices.md) - ⚠️ Must read
4. **For specific operations**: Reference the appropriate pattern guide below

## Core Patterns

### REST API Development
**Reference**: [references/rest-api-patterns/api-file-structure.md](references/rest-api-patterns/api-file-structure.md)

- API file structure with required settings (`go_package`, `group`, `compact_handler`)
- Three-layer architecture (Handler → Logic → Model)
- Request/response type definitions with validation
- Handler patterns and HTTP concerns
- Logic patterns and business implementation
- ✅ Correct vs ❌ incorrect patterns with examples

**When to use**: Creating or modifying REST API services, implementing HTTP endpoints

### Database Operations

- **[Best Practices](references/database-patterns/best-practices.md)**: Model import rules, error handling, field constants ⚠️
- **[Model Generation](references/database-patterns/model-generation.md)**: From SQL files or remote datasource
- **[Database Connection](references/database-patterns/database-connection.md)**: MySQL, PostgreSQL, SQLite, Redis configuration
- **[Condition Builder](references/database-patterns/condition-builder.md)**: Type-safe query building with `condition.NewChain()` API
- **[CRUD Operations](references/database-patterns/crud-operations.md)**: Generated methods (Insert, FindOne, Update, Delete, etc.)

**When to use**: Implementing data persistence, queries, or database operations

### Code Generation Workflow
**Reference**: [Project Structure](#project-structure)

jzero uses git-aware code generation - only changed files are regenerated

**Workflow**:
1. **Modify description file** (`desc/api/*.api`, `desc/sql/*.sql`, `desc/proto/*.proto`)
2. **Generate code**: `jzero gen --desc <file>`
3. **Implement business logic** in generated files

⚠️ Skipping step 2 causes compilation errors

**When to use**: Creating new features, modifying API definitions, adding database models

### Configuration Management
**Reference**: [Configuration](#configuration)

jzero supports flexible configuration with priority: `Environment Variables` > `CLI Flags` > `Config File`

- **CLI config** (`.jzero.yaml`): Code generation settings, git-change mode
- **App config** (`etc/etc.yaml`): REST, RPC, database, Redis settings
- **Environment overrides**: `export JZERO_GEN_GIT_CHANGE=true`

**When to use**: Setting up projects, configuring databases, adjusting generation behavior

## Project Structure

```
jzero-skills/
├── SKILL.md                           # This file - skill entry point
├── references/                        # Detailed pattern documentation
│   ├── rest-api-patterns/            # REST API guides
│   │   ├── README.md                 # Navigation index
│   │   └── api-file-structure.md     # ⚠️ Critical rules for .api files
│   └── database-patterns/            # Database operation guides
│       ├── README.md                 # Navigation index
│       ├── best-practices.md         # ⚠️ Critical rules with examples
│       ├── database-connection.md    # DB & Redis setup
│       ├── model-generation.md       # Generate models from SQL
│       ├── condition-builder.md      # Type-safe query building
│       └── crud-operations.md        # CRUD methods reference
```

**Typical jzero project structure**:
```
myproject/
├── desc/
│   ├── api/          # .api files → generates handlers
│   ├── sql/          # .sql files → generates models
│   └── proto/        # .proto files → generates RPC code
├── internal/
│   ├── handler/      # HTTP handlers (generated)
│   ├── logic/        # Business logic (implement here)
│   ├── model/        # Data models (generated)
│   ├── svc/          # Service context (dependencies)
│   ├── config/       # Config structs
│   └── middleware/   # Custom middleware
├── etc/
│   └── etc.yaml      # Configuration
└── .jzero.yaml       # jzero CLI config
```

## Common Workflows

### Creating a New REST API Endpoint

1. **Define API specification** in `.api` file with required settings:
   ```api
   info() { go_package: "user" }
   @server(group: user, compact_handler: true)
   ```
2. **Generate code**: `jzero gen --desc desc/api/user.api`
3. **Implement logic** in `internal/logic/` following three-layer architecture
4. **Test**: Use Swagger UI at `http://localhost:8000/swagger`

See detailed patterns: [REST API File Structure](references/rest-api-patterns/api-file-structure.md)

### Implementing Database Operations

1. **Create SQL schema** in `desc/sql/*.sql`
2. **Generate model**: `jzero gen --desc desc/sql/users.sql`
3. **Inject model** into ServiceContext
4. **Use condition builder** in logic layer
5. **Handle errors**

See detailed patterns: [Database Best Practices](references/database-patterns/best-practices.md)

### Setting Up Database Connection

1. **Configure in `etc/etc.yaml`**:
   ```yaml
   sqlx:
     driverName: mysql
     dataSource: "root:pass@tcp(127.0.0.1:3306)/mydb"
   redis:
     host: "127.0.0.1:6379"
     type: node
   ```
2. **Initialize in ServiceContext** with modelx.MustNewConn
3. **Register models** in Model struct

See detailed guide: [Database Connection](references/database-patterns/database-connection.md)

## Key Principles

### ✅ Always Follow

- **Three-layer architecture**: Handler → Logic → Model separation
- **API file requirements**: Set `go_package`, `group`, `compact_handler: true`
- **Condition builder**: Use `condition.NewChain()`, never `condition.New()`
- **Model imports**: Use alias `xxmodel "project/internal/model/xx"`
- **Error handling**: Use `errors.Is(err, model.ErrNotFound)` from `github.com/pkg/errors`
- **Code generation**: Run `jzero gen --desc` before implementing logic
- **Type safety**: Use generated field constants (e.g., `usersmodel.Id`)

### ❌ Never Do

- Put business logic in handlers (belongs in logic layer)
- Skip `go_package`, `group`, or `compact_handler` in `.api` files
- Use `condition.New()` instead of `condition.NewChain()`
- Import models without alias: `"project/internal/model/users"` (wrong)
- Use `==` for error comparison: `if err == ErrNotFound` (wrong)
- Hard-code configuration values
- Implement logic before generating code

## Progressive Learning

**New to jzero?**
1. Read this file (SKILL.md) for overview
2. Create a project: `jzero new myapi --frame api`
3. Follow [Code Generation Workflow](#code-generation-workflow)
4. Study [REST API File Structure](references/rest-api-patterns/api-file-structure.md)

**Building REST APIs?**
1. Master API file requirements (critical for avoiding regeneration issues)
2. Learn three-layer architecture patterns
3. Study [Database Best Practices](references/database-patterns/best-practices.md)
4. Reference [Condition Builder](references/database-patterns/condition-builder.md) for queries

**Working with databases?**
1. ⚠️ **Must read**: [Database Best Practices](references/database-patterns/best-practices.md)
2. Set up connection: [Database Connection](references/database-patterns/database-connection.md)
3. Generate models: [Model Generation](references/database-patterns/model-generation.md)
4. Learn CRUD operations: [CRUD Operations](references/database-patterns/crud-operations.md)

**Production deployment?**
1. Review all best practices in reference guides
2. Configure proper error handling and logging
3. Set up caching strategies with Redis
4. Implement resilience patterns (circuit breaker, rate limiting)

## Resources

- **Official documentation**: [docs.jzero.io](https://docs.jzero.io)
- **GitHub repository**: [jzero-io/jzero](https://github.com/jzero-io/jzero)
- **Examples**: [jzero-io/examples](https://github.com/jzero-io/examples)
- **Base framework**: [zeromicro/go-zero](https://github.com/zeromicro/go-zero)