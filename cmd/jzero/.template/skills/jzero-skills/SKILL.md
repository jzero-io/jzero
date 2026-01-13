---
name: jzero-skills
description: Comprehensive knowledge base for jzero framework (enhanced go-zero). Use this skill when working with jzero to understand correct patterns for REST APIs (Handler/Logic/Context architecture), RPC services (service discovery, load balancing), Gateway services, database operations (sqlx, MongoDB, caching), resilience patterns (circuit breaker, rate limiting), and jzero-specific features (git-change-based generation, flexible configuration, custom templates). Essential for generating production-ready jzero code that follows framework conventions.
license: Apache-2.0
---

# jzero Skills for AI Agents

Structured knowledge base optimized for AI agents to help developers work effectively with the [jzero](https://github.com/jzero-io/jzero) framework - an enhanced framework built on top of go-zero.

## Overview

This skill provides AI agents with comprehensive jzero knowledge to:
- Generate accurate code following jzero conventions
- Understand the three-layer architecture (Handler → Logic → Model)
- Apply jzero's enhanced features (simplify code generation, flexible configuration, custom templates)
- Work with API, RPC, and Gateway project types
- Build production-ready applications

## What is jzero?

jzero is an enhancement framework built on [go-zero](https://github.com/zeromicro/go-zero) and [goctl](https://github.com/zeromicro/go-zero/tree/master/tools/goctl) that provides:

- **Flexible configuration**: Control via YAML files, environment variables, and CLI flags
- **Simplify code generation**: Only generate code for changed files
- **Enhanced templating**: Custom templates with embedded `.tpl` support
- **Multiple project types**: API, RPC, and Gateway services
- **Serverless support**: Built-in serverless deployment capabilities
- **AI-friendly**: Optimized for AI-assisted development

## Quick Start

When helping with jzero development:

1. **For specific patterns**: Reference the appropriate pattern guide
2. **⚠️ CRITICAL for database operations**: Always consult [references/condition-builder.md](references/condition-builder.md) before writing any query conditions
3. **⚠️ CRITICAL for CRUD operations**: Follow best practices in [references/best-practices.md](references/best-practices.md) and [references/crud-operations.md](references/crud-operations.md)

## Core Patterns

### REST API Development
Reference: [references/rest-api-patterns.md](references/rest-api-patterns.md)

- Handler/Logic/Context three-layer architecture
- Request validation and error handling
- Middleware implementation (auth, logging, metrics)
- Response formatting with httpx
- Complete CRUD examples with ✅ correct and ❌ incorrect patterns

**When to use**: Creating or modifying REST API services, implementing HTTP endpoints

### Database Operations
Reference: [references/database-patterns.md](references/database-patterns.md)

- SQL operations with sqlx (CRUD, transactions, batch operations)
- Redis caching strategies
- Connection pooling and performance optimization
- Enhanced model generation with jzero

**When to use**: Implementing data persistence, caching, or database queries

### ⚠️ CRITICAL: Condition Package Usage

Reference: [references/condition-builder.md](references/condition-builder.md)

**‼️ MOST IMPORTANT RULE**: **ALWAYS use `condition.NewChain()` API for building query conditions**

The condition package provides two ways to build queries, but you MUST use the chain API:

```go
// ✅ CORRECT - ALWAYS use this pattern
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    usermodel "github.com/yourproject/internal/model/user"
)

conditions := condition.NewChain().
    Equal(usermodel.Username, req.Username).
    Greater(usermodel.Age, 18).
    Build()

// Use with any *ByCondition method
result, err := model.FindOneByCondition(ctx, nil, conditions...)

// ❌ WRONG - NEVER use this pattern
conditions := condition.New(
    condition.Condition{Field: usermodel.Username, Operator: condition.Equal, Value: req.Username},
)
```

**Key Points**:
- ✅ Use `condition.NewChain()` - Fluent, type-safe API
- ✅ Use generated field constants (e.g., `usermodel.Username`, `usermodel.Email`)
- ✅ Chain multiple conditions: `Equal().Like().In().Build()`

**When to use**: Writing ANY database query with conditions (filters, searches, pagination)

## jzero-Specific Features

### ⚠️ Critical: Code Generation Workflow

**IMPORTANT**: When modifying description files, you MUST run the corresponding `jzero gen` command BEFORE implementing any business logic:

```bash
# For API files (desc/api/*.api)
jzero gen --desc desc/api/your_file.api

# For Proto files (desc/proto/*.proto)
jzero gen --desc desc/proto/your_file.proto

# For SQL files (desc/sql/*.sql)
jzero gen --desc desc/sql/your_file.sql
```

**Why this is required**:
- Generates the Handler/Logic/Model skeleton code
- Creates type definitions and interfaces
- Ensures code structure matches the description
- Failing to do so will cause compilation errors

**Workflow**:
1. Modify `.api`/`.proto`/`.sql` file
2. Run the appropriate `jzero gen --desc` command
3. Implement business logic in the generated files
4. Test and commit

- Use `jzero gen` command when creating or updating api/sql/proto file

### Flexible Configuration System

jzero supports multi-source configuration with priority: **Environment Variables > Flags > Config File**

```yaml
# .jzero.yaml
gen:
  git-change: true

  # For subcommands
  zrpcclient:
    output: client
```

```bash
# Override with environment variables
export JZERO_GEN_GIT_CHANGE=true
export JZERO_GEN_ZRPCCLIENT_OUTPUT=client

# Or command-line flags
jzero gen zrpcclient --output client
```

**Configuration Priority**: `JZERO_GEN_GIT_CHANGE=true jzero gen` = `jzero gen --git-change=true`

### Project Types

jzero supports three project types:

```bash
# Create API project
jzero new myapi --frame api

# Create RPC project
jzero new myrpc --frame rpc

# Create Gateway project
jzero new mygateway --frame gateway
```

## Integration with AI Tools

jzero provides enhanced AI tooling:

- **[jzero-intellij](https://github.com/jzero-io/jzero-intellij)**: GoLand plugin for jzero
- **[jzero-admin](https://github.com/jzero-io/jzero-admin)**: Admin dashboard built with jzero
- **[jzero-action](https://github.com/marketplace/actions/jzero-action)**: GitHub Actions for jzero
- **[skills command](#)**: Built-in command to install Claude skills

## Project Structure

```
jzero-skills/
├── SKILL.md                    # This file - skill entry point
├── references/                 # Detailed pattern documentation
│   ├── rest-api-patterns.md    # REST API best practices
│   ├── database-patterns.md    # Database operation patterns
│   ├── condition-builder.md    # ⚠️ CRITICAL: Condition package usage
│   ├── best-practices.md       # Database best practices
│   ├── crud-operations.md      # CRUD operation examples
│   ├── model-generation.md     # Model generation guide
│   └── database-connection.md  # Database connection setup
```

## Common Workflows

### Creating a New REST API Service

1. Create project with `jzero new myproject --frame api`
2. Define API specification in `desc/api/` directory
3. Generate code with `jzero gen`
4. Implement business logic in `logic` layer
5. Add validation and error handling
6. Test with built-in swagger at `http://localhost:8001/swagger`

See complete workflow in [references/rest-api-patterns.md](references/rest-api-patterns.md)

### Implementing Database Operations

1. Design database schema
2. Add SQL file with `jzero add sql user.sql`
3. Generate model with `jzero gen`
4. Use sqlx for queries in logic layer
5. Handle transactions and errors properly

See complete patterns in [references/database-patterns.md](references/database-patterns.md)

### Adding Middleware

1. Create middleware function in `internal/middleware/`
2. Register in route configuration
3. Implement authentication/authorization logic
4. Pass data through request context
5. Handle errors appropriately

See middleware patterns in [references/rest-api-patterns.md](references/rest-api-patterns.md)

## Key Principles

### ✅ Always Follow

- **Three-layer architecture**: Handler → Logic → Model separation
- **Error handling**: Use structured errors, not `fmt.Errorf` in APIs
- **Configuration**: Use jzero's flexible configuration system
- **Context propagation**: Pass `ctx` through all layers
- **Type safety**: Define types in `.api` files, generate with jzero
- **Use jzero commands**: Prefer `jzero` over `goctl` for enhanced features
- **⚠️ Condition queries**: ALWAYS use `condition.NewChain()` API (see [condition-builder.md](references/condition-builder.md))
- **Field constants**: Use generated constants (e.g., `user.Username`) not hardcoded strings
- **References**: Consult [references/](references/) before implementing patterns

### ❌ Never Do

- Put business logic in handlers
- Ignore errors or use bare `fmt.Errorf` for HTTP errors
- Hard-code configuration values (use .jzero.yaml)
- Skip validation of user inputs
- Bypass the three-layer architecture
- Use `goctl` directly when jzero provides enhanced alternatives
- **‼️ Use `condition.New()` instead of `condition.NewChain()`** - This is critical!
- Use hardcoded strings for database field names
- Skip reading reference documentation before implementing database operations

## Progressive Learning

**New to jzero?**
1. Start with [getting-started/quick-start.md](getting-started/quick-start.md)
2. Build a simple REST API using [references/rest-api-patterns.md](references/rest-api-patterns.md)
3. Add database operations from [references/database-patterns.md](references/database-patterns.md)

## Resources

- **Official documentation**: [docs.jzero.io](https://docs.jzero.io)
- **GitHub repository**: [jzero-io/jzero](https://github.com/jzero-io/jzero)
- **Examples**: [jzero-io/examples](https://github.com/jzero-io/examples)
- **Community**: Join jzero developer community (see main repo README)

## Version Compatibility

This skill targets jzero 1.1+. jzero maintains compatibility with go-zero 1.5+. Patterns are updated regularly to reflect framework evolution. Always check official documentation for the latest API changes.

## Key Differences from go-zero

While jzero maintains full compatibility with go-zero, it adds:

1. **Enhanced CLI**: `jzero` command with more features than `goctl`
2. **Flexible configuration**: Multi-source config (YAML + ENV + CLI)
3. **Git-aware generation**: Only generate changed files
4. **Custom templates**: Built-in template customization
5. **Gateway support**: API Gateway project type
6. **Serverless**: Built-in serverless deployment
7. **Better DX**: Improved error messages, progress tracking, etc.

When in doubt, prefer jzero commands over goctl for the best experience.
