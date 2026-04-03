---
title: "AI-Driven Go Development: How jzero Makes Development More Efficient and Reliable"
icon: /icons/streamline-ultimate-blog-blogger-logo.svg
---

As AI-assisted programming gradually becomes mainstream, everyone is gradually using AI tools to assist development, which has indeed greatly improved coding efficiency. However, an overlooked issue is:

**AI-generated code often lacks readability and maintainability, making it difficult to integrate into team development standards and engineering systems.**

### Pain Points of Traditional AI Programming

**1. "Tug of war" with AI consumes a lot of energy**

When using traditional AI tools, you need to repeatedly "tug" with AI to get usable code:

```
You: Help me write a user registration interface
AI: Generated a piece of code (doesn't follow framework standards)

You: Naming doesn't follow team standards, refer to user/list.go
AI: Modified another version (still not completely right)

You: Should use the framework's Model layer methods, don't write SQL directly
AI: Continue modifying...

You: Compilation error, fix the problem
AI: Continue modifying...

You: Finally works (but spent 30 minutes, experienced n rounds of dialogue)
```

**2. AI is only goal-oriented, ignoring engineering practices**

- ❌ **Lacks best practices**: AI doesn't know framework-recommended patterns, generated code "works but isn't elegant"
- ❌ **Doesn't follow team standards**: Naming, comments, directory structure vary
- ❌ **Prone to repeated modifications**: Same requirement, different code generated each time
- ❌ **Uncontrollable quality**: Lacks error handling, logging, monitoring and other engineering elements

**3. Cannot reference framework best practices**

- ❌ Doesn't know what utility classes and helper functions the framework provides
- ❌ Doesn't understand framework-recommended directory structure and module division
- ❌ Unclear about framework performance optimization suggestions
- ❌ Generated code looks "foreign", difficult to integrate into existing projects

### How jzero Solves These Problems

**By defining workflows and best practices in agent skills, reduce the "tug of war" between users and AI**:

1. **Standardized workflow**: `api/proto/sql` → `jzero gen` → framework code → AI fills business logic
2. **Complete framework knowledge**: jzero-skills contains framework best practices, AI follows automatically
3. **Work within constraints**: AI intelligently fills business logic under framework constraints, both efficient and standardized
4. **One-time generation is usable**: Reduce repeated modifications, guaranteed code quality

- 🔧 **Engine One** (generate basic framework code based on fixed templates): 100% standardized, ensures architectural consistency
- 🤖 **Engine Two** (AI Skills generate business code): Intelligently fill business logic under framework constraints
- 💥 **Combined**: Enjoy AI efficiency while ensuring code quality and maintainability

## jzero's Dual-Engine Development Mode

jzero's dual-engine mode divides code generation into two stages, where AI first generates descriptor files, then generates framework code and business logic based on descriptor files:

**🏗️ jzero Dual-Engine Architecture**

---

**🔧 Engine One: Basic Framework Code Generation**

*   AI generates api file → jzero framework generates Handler/Logic/Types
*   AI generates proto file → jzero framework generates Server/Logic/Pb
*   AI generates sql file → jzero generates Model layer general method code
*   **Features**: Standardized, predictable, follows best practices
*   **Speed**: Second-level generation

---

**🤖 Engine Two: AI Skills Business Code Generation**

*   Logic layer business logic implementation
*   Complex queries and transaction processing
*   **Features**: Intelligent, automated, understands business requirements
*   **Speed**: Minute-level implementation

---

**💥 Combined**: Complete functional code, 10x development efficiency improvement

### Understanding Dual-Engine in One Minute

**Traditional Development** (you do everything yourself):
```
Requirement → Hand-write .api → Hand-write Handler → Hand-write Logic → Hand-write Model → Debug
       20min     10min      60min      30min     30min
Total: 2.5 hours
```

**jzero Dual-Engine** (you only describe requirements):
```
Requirement → AI creates .api → jzero framework generates basic framework → AI generates business logic
       1min         10sec            3-5min
Total: 7 minutes ⚡
```

**Key Difference**:
- 🔧 **Engine One** (jzero gen): Generate standardized framework code, ensure architectural consistency
- 🤖 **Engine Two** (AI Skills): Understand requirements, intelligently fill business logic
- 💡 **Your role**: From "code porter" to "architect and business expert"

### Why Do We Need This Mode?

**Pain points of traditional development mode**:

- ❌ **Reinventing the wheel**: Every new feature requires writing similar Handler/Logic/Model
- ❌ **Inconsistent architecture**: Different developers write code with different styles, high maintenance cost
- ❌ **AI doesn't understand frameworks**: Ordinary AI-generated code needs lots of adjustments to integrate into projects
- ❌ **Efficiency bottleneck**: Framework code takes a lot of time, little time left for business logic

**Advantages of jzero dual-engine mode**:

| Dimension | Engine One: Framework Generation | Engine Two: AI Skills | Combined Advantage    |
|------|----------------|------------------|---------|
| **Generated Content** | Handler/Logic/Model framework | Business logic implementation | Complete functional code  |
| **Code Quality** | 100% follows framework standards | Follows best practices | Production level    |
| **Development Speed** | Second-level generation | Minute-level implementation | 10x improvement  |
| **Architectural Consistency** | ✅ Completely consistent | ✅ Automatically follows | Frictionless team collaboration |
| **Maintainability** | Standardized structure | Clear business logic | Extremely easy to maintain    |

---

## Dual-Engine in Practice: How to Collaborate to Generate High-Quality, Maintainable Code

Let's see how two engines work together perfectly through a specific example.

Your requirement:
```
jzero-skills create user management api, supporting:
1. User registration (username 3-20 characters, email validation, password at least 8 characters)
2. Get user info
```

### 🔧 Engine One: Basic Framework Code Generation (Automatic)

**Step 1: AI creates api definition file**

Based on requirement description, AI automatically infers and adds validation rules:

- ✅ Username 3-20 characters → `validate:"required,min=3,max=20"`
- ✅ Email validation → `validate:"required,email"`
- ✅ Password at least 8 characters → `validate:"required,min=8"`

AI automatically generates `desc/api/user.api`:

```go
info(
    go_package: "user"
)

type User {
    id       int64  `json:"id"`
    username string `json:"username"`
    email    string `json:"email"`
}

type RegisterRequest {
    username string `json:"username" validate:"required,min=3,max=20"`
    email    string `json:"email" validate:"required,email"`
    password string `json:"password" validate:"required,min=8"`
}

type RegisterResponse {
	User
}

type GetRequest {
	id int64 `path:"id"`
}

type GetResponse {
	User
}

@server(
    prefix: /api/user
    group: user
    jwt: JwtAuth
    middleware: AuthX
    compact_handler: true
)
service simpleapi {
    @doc "User registration"
    @handler Register
    post /register (RegisterRequest) returns (RegisterResponse)

    @doc "Get user info"
    @handler Get
    get /:id (GetRequest) returns (GetResponse)
}
```

**Step 2: AI automatically creates database table definition**

AI automatically generates `desc/sql/user.sql` based on requirements:

```sql
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'User ID',
  `username` varchar(20) NOT NULL COMMENT 'Username',
  `email` varchar(255) NOT NULL COMMENT 'Email',
  `password` varchar(255) NOT NULL COMMENT 'Password',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User table';
```

AI intelligently designs:
- ✅ Automatically infers field types based on validation rules (username 20 characters → varchar(20))
- ✅ Automatically adds indexes (username/email unique indexes)
- ✅ Automatically adds timestamp fields (created_at/updated_at)
- ✅ Automatically adds table and field comments

**Step 2: AI automatically executes framework generation command**

```bash
# 1. Generate Handler/Logic/Types (from api file)
jzero gen --desc desc/api/user.api

# 2. Generate Model layer (from sql file)
jzero gen --desc desc/sql/user.sql
```

→ Generated files:
- ✅ `internal/model/model.go` // All model registration
- ✅ `internal/model/user/*.go` // User Model layer (CRUD methods, field constants)
- ✅ `internal/handler/user/user_compact.go` // user handler
- ✅ `internal/handler/routes.go` // Route registration
- ✅ `internal/logic/user/register.go` // Register business logic empty implementation
- ✅ `internal/logic/user/get.go` // Get business logic empty implementation
- ✅ `internal/types/user/types.go`   // type struct definition

### 🤖 Engine Two: AI Skills Fill Business Logic (AI Intelligent Generation)

**Step 3: AI implements register and get business logic**

**Register framework code generated by Engine One** `internal/logic/user/register.go`:

```go
package user

import (
    "context"
    "net/http"

    "github.com/zeromicro/go-zero/core/logx"

    "simpleapi/internal/svc"
    types "simpleapi/internal/types/user"
)

type Register struct {
    logx.Logger
    ctx     context.Context
    svcCtx  *svc.ServiceContext
    r       *http.Request
}

func NewRegister(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *Register {
    return &Register{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
        r:      r,
    }
}

func (l *Register) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
    // todo: add your logic here and delete this line

    return
}
```

**After Engine Two AI Skills fills register business logic**:

```go
package user

import (
    "context"
    "net/http"

    "github.com/pkg/errors"
    "github.com/jzero-io/jzero/core/stores/condition"
    usermodel "github.com/yourproject/internal/model/user"
    "github.com/zeromicro/go-zero/core/logx"
    "golang.org/x/crypto/bcrypt"

    "simpleapi/internal/svc"
    types "simpleapi/internal/types/user"
)

type Register struct {
    logx.Logger
    ctx     context.Context
    svcCtx  *svc.ServiceContext
    r       *http.Request
}

func NewRegister(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *Register {
    return &Register{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
        r:      r,
    }
}

func (l *Register) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
    chain := condition.NewChain().
        Equal(usermodel.Username, req.Username)

    existingUser, err := l.svcCtx.Model.User.FindOneByCondition(l.ctx, nil, chain.Build()...)
    if err != nil {
        if !errors.Is(err, usermodel.ErrNotFound) {
            l.Logger.Errorf("failed to check username existence: %v", err)
            return nil, err
        }
    } else {
        return nil, errors.New("username already exists")
    }

    user := &usermodel.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
    }

    err = l.svcCtx.Model.User.InsertV2(l.ctx, nil, user)
    if err != nil {
        l.Logger.Errorf("failed to insert user: %v", err)
        return nil, err
    }

    return &types.RegisterResponse{
        Id:       user.Id,
        Username: user.Username,
        Email:    user.Email,
    }, nil
}
```

**Get framework code generated by Engine One** `internal/logic/user/get.go`:

```go
package user

import (
    "context"
    "net/http"

    "github.com/zeromicro/go-zero/core/logx"

    "simpleapi/internal/svc"
    types "simpleapi/internal/types/user"
)

type Get struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
    r      *http.Request
}

// Get user info
func NewGet(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *Get {
    return &Get{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
        r:      r,
    }
}

func (l *Get) Get(req *types.GetRequest) (resp *types.GetResponse, err error) {
    // todo: add your logic here and delete this line

    return
}
```

**After Engine Two AI Skills fills get business logic**:

```go
package user

import (
    "context"
    "net/http"

    "github.com/pkg/errors"
    usermodel "github.com/yourproject/internal/model/user"
    "github.com/zeromicro/go-zero/core/logx"

    "simpleapi/internal/svc"
    types "simpleapi/internal/types/user"
)

type Get struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
    r      *http.Request
}

// Get user info
func NewGet(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *Get {
    return &Get{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
        r:      r,
    }
}

func (l *Get) Get(req *types.GetRequest) (resp *types.GetResponse, err error) {
    user, err := l.svcCtx.Model.User.FindOne(l.ctx, nil, req.Id)
    if err != nil {
        if errors.Is(err, usermodel.ErrNotFound) {
            return nil, errors.New("user does not exist")
        }
        l.Logger.Errorf("failed to find user: %v", err)
        return nil, err
    }

    return &types.GetResponse{
        Id:       user.Id,
        Username: user.Username,
        Email:    user.Email,
    }, nil
}
```

### 📊 Summary: Core Value of Dual-Engine Collaboration

Through this complete example, we can clearly see the workflow and value of jzero dual-engine mode collaboration:


| Stage | Responsible Engine | Core Task | Output |
|------|---------|---------|---------|
| **Descriptor File Generation** | 🤖 AI Skills | Understand requirements, generate api/sql definitions | Standardized descriptor files |
| **Framework Code Generation** | 🔧 jzero gen | Generate standard framework based on descriptor files | Handler/Logic/Model empty implementations |
| **Business Logic Filling** | 🤖 AI Skills | Implement specific business logic | Complete usable functional code |

**Key Advantages**

✅ **Each specializes in their field**
   - Engine One (jzero gen): Ensures architectural consistency, code standardization, fast generation
   - Engine Two (AI Skills): Understands business, intelligently fills logic, handles complex logic

✅ **Quality and efficiency both considered**
   - Framework code 100% follows standards, no technical debt
   - Business logic intelligently generated, reduce repetitive work
   - 10x development efficiency improvement, code quality not compromised but improved

✅ **Predictable output**
   - Framework code location, naming, structure completely standardized
   - Team members generate code with consistent style
   - Easy for Code Review and team collaboration

✅ **Reduce cognitive burden**
   - Developers only focus on business logic, framework details automatically handled
   - New team members get up to speed quickly, low project maintenance cost
   - AI works under constraints, won't "randomly create"

**Comparison with Traditional Development**

| Development Mode | Time | Code Quality | Architectural Consistency | Maintainability |
|---------|------|---------|---------|---------|
| **Pure Manual Development** | 2.5 hours | ⭐⭐⭐ | ❌ Uncertain | ⭐⭐ |
| **Traditional AI Assistance** | 30 minutes+ | ⭐⭐ | ❌ Uncertain | ⭐⭐ |
| **jzero Dual-Engine** | 7 minutes | ⭐⭐⭐⭐⭐ | ✅ Completely consistent | ⭐⭐⭐⭐⭐ |

**Core Philosophy**

jzero dual-engine mode embodies our core view on AI-assisted development:

> **AI should not replace developers in making architectural decisions, but should intelligently fill business logic under framework constraints.**


### Reflection: Why jzero Only Uses Skills and Abandons MCP

In exploring AI-assisted development, jzero attempted to integrate MCP (Model Context Protocol). This is a protocol that allows AI models to access external tools and data through standardized protocols. However, after in-depth practice and validation, we finally decided to abandon MCP.

**The reason is simple: jzero's code generation process is inherently simple and direct.**

**jzero's Core Design Philosophy**

jzero has adopted an extremely simple code generation approach from the start:

```
Descriptor files (api/sql/proto) → jzero gen → Complete framework code
```

The core advantages of this process:

✅ **File locations and structures are all conventional**
- api descriptor files uniformly placed in `desc/api/` directory
- sql descriptor files uniformly placed in `desc/sql/` directory
- proto descriptor files uniformly placed in `desc/proto/` directory

✅ **One command handles code generation**
- `jzero gen` // Automatically scan api/proto/sql to generate code
- `jzero gen --desc desc/api/user.api` // Only generate user.api
- `jzero gen --desc desc/proto/user.proto` // Only generate user.proto
- `jzero gen --desc desc/sql/user.sql` // Only generate user.sql

For a framework like jzero with clear conventions and simple commands, **skills is more direct, simple, efficient, and easy to use than mcp**.

## Quick Start

```bash
# 1. Install jzero
go install github.com/jzero-io/jzero/cmd/jzero@latest

# 2. Initialize jzero-skills
jzero skills init

# 3. Create project
jzero new simpleapi --frame api

# 4. Describe requirements in Claude Code
# "Use jzero-skills to create a user management api..."
```

**Find it useful? Please give jzero a ⭐ Star to support our continued improvement!**

GitHub: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)
