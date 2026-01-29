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
- Use jzero-specific features
- Build production-ready applications

## Quick Start

When helping with jzero development:

1. **For new projects**: Start with [Development Workflows](#development-workflows)
2. **For REST APIs**: Check [REST API File Structure](references/rest-api-patterns/api-file-structure.md) - ⚠️ Critical rules
3. **For RPC services**: Review [Proto File Structure](references/rpc-patterns/proto-file-structure.md) - Proto standards & multi-proto support
4. **For databases**: Review [Database Best Practices](references/database-patterns/best-practices.md) - ⚠️ Must read
5. **For SQL changes**: Check [SQL Migration Guide](references/database-patterns/sql-migration.md) - ⚠️ Schema changes
6. **For specific operations**: Reference the appropriate pattern guide below

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

### RPC Services

- **[Proto File Structure](references/rpc-patterns/proto-file-structure.md)**: Proto standards, multi-proto support, file structure, HTTP gateway, OpenAPI docs
- **[Proto Field Validation](references/rpc-patterns/proto-validation.md)**: Field validation with protovalidate, CEL expressions, built-in constraints
- **[Proto Middleware](references/rpc-patterns/proto-middleware.md)**: HTTP/RPC middleware at service and method levels

**When to use**: Creating or modifying RPC services, working with proto files, adding validation or middleware

### Database Operations

- **[Best Practices](references/database-patterns/best-practices.md)**: Model import rules, condition chain usage, error handling, field constants ⚠️🚨
- **[SQL Migration Guide](references/database-patterns/sql-migration.md)**: Managing schema changes with up/down migrations ⚠️
- **[Model Generation](references/database-patterns/model-generation.md)**: From SQL files or remote datasource
- **[Database Connection](references/database-patterns/database-connection.md)**: MySQL, PostgreSQL, SQLite, Redis configuration
- **[CRUD Operations](references/database-patterns/crud-operations.md)**: Generated methods (Insert, FindOne, Update, Delete, etc.)

**⚠️ CRITICAL REMINDER**: ALWAYS use `condition.NewChain()` - NEVER use `condition.New()`

**When to use**: Implementing data persistence, queries, or database operations

### Development Workflows
**Reference**: [Project Structure](#project-structure)

#### Creating a New REST API Endpoint

1. **Define API specification** in `.api` file with required settings:
2. **Generate api code**: `jzero gen --desc desc/api/user.api`
3. **Implement logic** in `internal/logic/` following three-layer architecture

See detailed patterns: [REST API File Structure](references/rest-api-patterns/api-file-structure.md)

#### Implementing Database Operations

**Choose your schema mode first:**

**Local SQL Mode** (schema files in `desc/sql/`):
1. **Create/update SQL schema** in `desc/sql/*.sql`
2. **Create migration files** in `desc/sql_migration/` (xx.up.sql & xx.down.sql) ⚠️
3. **Apply migrations** (development: `jzero migrate up`, production: auto in `cmd/server.go`)
4. **Generate model**: `jzero gen --desc desc/sql/users.sql`

**Remote Datasource Mode** (schema from live database):
1. **Create migration files** in `desc/sql_migration/` (xx.up.sql & xx.down.sql) ⚠️
2. **Apply migrations** (development: `jzero migrate up`, production: auto in `cmd/server.go`)
3. **Generate model**: `jzero gen`

**Common steps (both modes)**:
- Inject model into ServiceContext
- Use condition builder in logic layer
- Handle errors properly

⚠️ **Migration rules**: Always create both up/down files, use consecutive numbering (1, 2, 3...)

See detailed patterns: [SQL Migration Guide](references/database-patterns/sql-migration.md) | [Database Best Practices](references/database-patterns/best-practices.md)

#### Setting Up Database Connection

1. **Configure in `etc/etc.yaml`**:
2. **Initialize in ServiceContext** with modelx.MustNewConn
3. **Register models** in Model struct

See detailed guide: [Database Connection](references/database-patterns/database-connection.md)

## Project Structure

```
jzero-skills/
├── SKILL.md                           # This file - skill entry point
├── references/                        # Detailed pattern documentation
│   ├── rest-api-patterns/            # REST API guides
│   │   └── api-file-structure.md     # ⚠️ Critical rules for .api files
│   ├── rpc-patterns/                # RPC/Proto service guides
│   │   ├── proto-file-structure.md  # Proto standards & multi-proto
│   │   ├── proto-validation.md      # Field validation guide
│   │   └── proto-middleware.md      # Middleware patterns
│   └── database-patterns/            # Database operation guides
│       ├── best-practices.md         # ⚠️ Critical rules with examples
│       ├── sql-migration.md          # ⚠️ Schema changes & migrations
│       ├── database-connection.md    # DB & Redis setup
│       ├── model-generation.md       # Generate models from SQL
│       └── crud-operations.md        # CRUD methods reference
```

**Typical jzero project structure**:
```
myproject/
├── .jzero.yaml       # CLI config: code generation, ⚠️ migrate settings
├── desc/
│   ├── api/          # .api files → generates handlers
│   ├── sql/          # .sql files → generates models (local SQL mode)
│   ├── sql_migration/ # xx.up.sql & xx.down.sql for schema changes ⚠️
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

## Key Principles

### ✅ Always Follow

- **🚨 Condition builder**: ALWAYS use `condition.NewChain()`, NEVER use `condition.New()` - **THIS IS CRITICAL** 🚨
- **Three-layer architecture**: Handler → Logic → Model separation
- **API file requirements**: Set `go_package`, `group`, `compact_handler: true`
- **Model imports**: Use alias `xxmodel "project/internal/model/xx"`
- **Error handling**: Use `errors.Is(err, model.ErrNotFound)` from `github.com/pkg/errors`
- **Code generation**: Run `jzero gen --desc` before implementing logic
- **Type safety**: Use generated field constants (e.g., `usersmodel.Id`)

### ❌ Never Do

- 🚫 **NEVER use `condition.New()`** - This is error-prone and deprecated. **ALWAYS use `condition.NewChain()`**
- Put business logic in handlers (belongs in logic layer)
- Skip `go_package`, `group`, or `compact_handler` in `.api` files
- Import models without alias: `"project/internal/model/users"` (wrong)
- Use `==` for error comparison: `if err == ErrNotFound` (wrong)
- Hard-code configuration values
- Implement logic before generating code

## Progressive Learning

**New to jzero?**
1. Read this file (SKILL.md) for overview
2. Create a project: `jzero new myapi --frame api`
3. Follow [Development Workflows](#development-workflows)
4. Study [REST API File Structure](references/rest-api-patterns/api-file-structure.md)

**Building REST APIs?**
1. Master API file requirements (critical for avoiding regeneration issues)
2. Learn three-layer architecture patterns
3. Study [Database Best Practices](references/database-patterns/best-practices.md)

**Creating RPC services?**
1. Learn proto file structure: [Proto File Structure](references/rpc-patterns/proto-file-structure.md)
2. Add validation: [Proto Field Validation](references/rpc-patterns/proto-validation.md)
3. Implement middleware: [Proto Middleware](references/rpc-patterns/proto-middleware.md)

**Working with databases?**
1. ⚠️ **Must read**: [Database Best Practices](references/database-patterns/best-practices.md)
2. ⚠️ **Must read**: [SQL Migration Guide](references/database-patterns/sql-migration.md) - Schema changes
3. Set up connection: [Database Connection](references/database-patterns/database-connection.md)
4. Generate models: [Model Generation](references/database-patterns/model-generation.md)
5. Learn CRUD operations: [CRUD Operations](references/database-patterns/crud-operations.md)

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