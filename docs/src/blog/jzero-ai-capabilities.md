---
title: AI 驱动的 Go 开发：jzero 如何让开发更高效更可靠
icon: streamline-ultimate:blog-blogger-logo
---

在 AI 辅助编程逐渐成为主流的今天，大家都逐渐在用 AI 工具辅助开发了，也确实大大提升了编码效率。但一个不可忽视的问题是：

**AI 生成的代码往往缺乏可读性与可维护性，难以融入团队的开发标准和工程体系。**

### 传统 AI 编程的痛点

**1. 与 AI 的"拉扯"消耗大量精力**

使用传统 AI 工具时，你需要反复与 AI "拉扯"才能获得可用的代码：

```
你：帮我写一个用户注册接口
AI：生成了一段代码（不符合框架规范）

你：命名不符合团队规范，参照 user/list.go 的写法
AI：又改了一版（还是不完全对）

你：要使用框架的 Model 层方法，不要直接写 SQL
AI：继续修改...

你：编译错误，修复问题
AI：继续修改...

你：终于能用了（但已经花了 30 分钟，经历了 n 轮对话）
```

**2. AI 仅以目标为导向，忽略工程实践**

- ❌ **缺乏最佳实践**：AI 不知道框架推荐写法，生成的代码"能用但不优雅"
- ❌ **不符合团队规范**：命名、注释、目录结构各不相同
- ❌ **易反复修改**：同样的需求，每次生成的代码都不一样
- ❌ **质量不可控**：缺乏错误处理、日志、监控等工程要素

**3. 无法参照框架最佳实践**

- ❌ 不知道框架提供了哪些工具类和辅助函数
- ❌ 不了解框架推荐的目录结构和模块划分
- ❌ 不清楚框架的性能优化建议
- ❌ 生成的代码像"外来的"，难以融入现有项目

### jzero 如何解决这些问题

**通过 agent skills 中定义的工作流程以及最佳实践，减少用户与 AI 之间的"拉扯"**：

1. **标准化工作流**：`api/proto/sql` → `jzero gen` → 框架代码 → AI 填充业务逻辑
2. **完整的框架知识**：jzero-skills 包含框架最佳实践，AI 自动遵循
3. **约束条件下发挥**：AI 在框架约束下智能填充，既高效又规范
4. **一次生成即可用**：减少反复修改，代码质量有保障

- 🔧 **引擎一**（基于固定模板生成基础框架代码）：100% 标准化，保证架构统一
- 🤖 **引擎二**（AI Skills 生成业务代码）：在框架约束下智能填充业务逻辑
- 💥 **双剑合璧**：既享受 AI 的效率，又保证代码质量和可维护性

## jzero 的双引擎开发模式

jzero 双引擎模式将代码生成分为两个阶段，AI 先生成描述文件，再基于描述文件生成框架代码和业务逻辑：

**🏗️ jzero 双引擎架构**

---

**🔧 引擎一：基础框架代码生成**

*   AI 生成 api 文件 → jzero 框架生成 Handler/Logic/Types
*   AI 生成 proto 文件 → jzero 框架生成 Server/Logic/Pb
*   AI 生成 sql 文件 → jzero 生成 Model 层通用方法代码
*   **特点**：标准化、可预测、符合最佳实践
*   **速度**：秒级生成

---

**🤖 引擎二：AI Skills 业务代码生成**

*   Logic 层业务逻辑实现
*   复杂查询和事务处理
*   **特点**：智能化、自动化、理解业务需求
*   **速度**：分钟级实现

---

**💥 双剑合璧**：完整功能代码，开发效率提升 10 倍

### 一分钟看懂双引擎

**传统开发**（你自己做所有事）：
```
需求 → 手写 .api → 手写 Handler → 手写 Logic → 手写 Model → 调试
       20分钟     10分钟      60分钟      30分钟     30分钟
总计：2.5 小时
```

**jzero 双引擎**（你只描述需求）：
```
需求 → AI 创建 .api → jzero 框架生成基础框架 → AI 生成业务逻辑
       1分钟         10秒            3-5分钟
总计：7 分钟 ⚡
```

**关键区别**：
- 🔧 **引擎一**（jzero gen）：生成标准化框架代码，保证架构统一
- 🤖 **引擎二**（AI Skills）：理解需求，智能填充业务逻辑
- 💡 **你的角色**：从"代码搬运工"变成"架构师和业务专家"

### 为什么需要这种模式？

**传统开发模式的痛点**：

- ❌ **重复造轮子**：每次新建功能都要写一遍相似的 Handler/Logic/Model
- ❌ **架构不一致**：不同开发者写的代码风格各异，维护成本高
- ❌ **AI 不懂框架**：普通 AI 生成的代码需要大量调整才能融入项目
- ❌ **效率瓶颈**：框架代码占用大量时间，真正用于业务逻辑的时间少

**jzero 双引擎模式的优势**：

| 维度 | 引擎一：框架生成 | 引擎二：AI Skills | 结合优势    |
|------|----------------|------------------|---------|
| **生成内容** | Handler/Logic/Model 框架 | 业务逻辑实现 | 完整功能代码  |
| **代码质量** | 100% 符合框架规范 | 遵循最佳实践 | 生产级别    |
| **开发速度** | 秒级生成 | 分钟级实现 | 10 倍提升  |
| **架构统一** | ✅ 完全统一 | ✅ 自动遵循 | 团队协作无摩擦 |
| **可维护性** | 标准化结构 | 清晰的业务逻辑 | 极易维护    |

---

## 双引擎实战：如何协作生成高质量易维护的代码

让我们通过一个具体例子，看看两个引擎如何完美配合。

你的需求：
```
jzero-skills 创建用户管理 api，支持：
1. 用户注册（用户名 3-20 位、邮箱验证、密码至少 8 位）
2. 获取用户信息
```

### 🔧 引擎一：基础框架代码生成（自动执行）

**步骤 1：AI 创建 api 定义文件**

AI 根据需求描述，自动推断并添加验证规则：

- ✅ 用户名 3-20 位 → `validate:"required,min=3,max=20"`
- ✅ 邮箱验证 → `validate:"required,email"`
- ✅ 密码至少 8 位 → `validate:"required,min=8"`

AI 自动生成 `desc/api/user.api`：

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
    @doc "用户注册"
    @handler Register
    post /register (RegisterRequest) returns (RegisterResponse)

    @doc "获取用户信息"
    @handler Get
    get /:id (GetRequest) returns (GetResponse)
}
```

**步骤 2：AI 自动创建数据库表定义**

AI 根据需求自动生成 `desc/sql/user.sql`：

```sql
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `email` varchar(255) NOT NULL COMMENT '邮箱',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

AI 智能设计：
- ✅ 根据验证规则自动推断字段类型（username 20 位 → varchar(20)）
- ✅ 自动添加索引（username/email 唯一索引）
- ✅ 自动添加时间戳字段（created_at/updated_at）
- ✅ 自动添加表注释和字段注释

**步骤 2：AI 自动执行框架生成命令**

```bash
# 1. 生成 Handler/Logic/Types（从 api 文件）
jzero gen --desc desc/api/user.api

# 2. 生成 Model 层（从 sql 文件）
jzero gen --desc desc/sql/user.sql
```

→ 生成文件：
- ✅ `internal/model/model.go` // 所有 model 注册
- ✅ `internal/model/user/*.go` // User Model 层（CRUD 方法、字段常量）
- ✅ `internal/handler/user/user_compact.go` // user handler
- ✅ `internal/handler/routes.go` // 路由注册
- ✅ `internal/logic/user/register.go` // register 业务逻辑空实现
- ✅ `internal/logic/user/get.go` //  get 业务逻辑空实现
- ✅ `internal/types/user/types.go`   // type struct 定义

### 🤖 引擎二：AI Skills 填充业务逻辑（AI 智能生成）

**步骤 3：AI 实现 register 和 get 业务逻辑**

**引擎一生成的 register 框架代码** `internal/logic/user/register.go`：

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

**引擎二 AI Skills 填充 register 业务逻辑后**：

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
        return nil, errors.New("用户名已存在")
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

**引擎一生成的 get 框架代码** `internal/logic/user/get.go`：

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

// 获取用户信息
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

**引擎二 AI Skills 填充 get 业务逻辑后**：

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

// 获取用户信息
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
            return nil, errors.New("用户不存在")
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

### 📊 总结：双引擎协作的核心价值

通过这个完整的例子，我们可以清楚地看到 jzero 双引擎模式的协作流程和价值：


| 阶段 | 负责引擎 | 核心任务 | 输出成果 |
|------|---------|---------|---------|
| **描述文件生成** | 🤖 AI Skills | 理解需求，生成 api/sql 定义 | 规范的描述文件 |
| **框架代码生成** | 🔧 jzero gen | 基于描述文件生成标准框架 | Handler/Logic/Model 空实现 |
| **业务逻辑填充** | 🤖 AI Skills | 实现具体的业务逻辑 | 完整可用的功能代码 |

**关键优势**

✅ **各司其职，发挥所长**
   - 引擎一（jzero gen）：确保架构统一、代码规范、生成快速
   - 引擎二（AI Skills）：理解业务、智能填充、处理复杂逻辑

✅ **质量与效率兼顾**
   - 框架代码 100% 符合规范，无技术债
   - 业务逻辑智能生成，减少重复劳动
   - 开发效率提升 10 倍，代码质量不降反升

✅ **可预测的输出**
   - 框架代码位置、命名、结构完全标准化
   - 团队成员生成的代码风格一致
   - 便于 Code Review 和团队协作

✅ **降低认知负担**
   - 开发者只需关注业务逻辑，框架细节自动处理
   - 新人上手快，项目维护成本低
   - AI 在约束条件下工作，不会"胡乱发挥"

**与传统开发的对比**

| 开发模式 | 耗时 | 代码质量 | 架构统一 | 可维护性 |
|---------|------|---------|---------|---------|
| **纯手工开发** | 2.5 小时 | ⭐⭐⭐ | ❌ 不确定 | ⭐⭐ |
| **传统 AI 辅助** | 30 分钟+ | ⭐⭐ | ❌ 不确定 | ⭐⭐ |
| **jzero 双引擎** | 7 分钟 | ⭐⭐⭐⭐⭐ | ✅ 完全统一 | ⭐⭐⭐⭐⭐ |

**核心理念**

jzero 双引擎模式体现了我们对 AI 辅助开发的核心观点：

> **AI 不应该替代开发者做架构决策，而应该在框架约束下智能填充业务逻辑。**


### 思考: jzero 为何仅使用 skills 放弃 MCP

在探索 AI 辅助开发的过程中，jzero 曾经尝试过集成 MCP (Model Context Protocol)。这是一个允许 AI 模型通过标准化协议访问外部工具和数据的协议。然而，经过深入实践和验证，我们最终决定放弃 MCP。

**原因很简单：jzero 的代码生成流程本身就足够简单直接。**

**jzero 的核心设计理念**

jzero 从一开始就采用了一种极其简单的代码生成方式：

```
描述文件 (api/sql/proto) → jzero gen → 完整的框架代码
```

这个流程的核心优势在于：

✅ **文件位置和结构都是约定的**
- api 描述文件统一放在 `desc/api/` 目录
- sql 描述文件统一放在 `desc/sql/` 目录
- proto 描述文件统一放在 `desc/proto/` 目录

✅ **一条命令搞定代码生成**
- `jzero gen` // 自动扫描 api/proto/sql 生成代码
- `jzero gen --desc desc/api/user.api` // 仅生成 user.api
- `jzero gen --desc desc/proto/user.proto` // 仅生成 user.proto
- `jzero gen --desc desc/sql/user.sql` // 仅生成 user.sql

对于 jzero 这样约定明确、命令简单的框架，**skills 比 mcp 更直接、更简单高效、更易用**。

## 快速开始

```bash
# 1. 安装 jzero
go install github.com/jzero-io/jzero/cmd/jzero@latest

# 2. 初始化 jzero-skills
jzero skills init

# 3. 创建项目
jzero new simpleapi --frame api

# 4. 在 Claude Code 中描述需求
# "使用 jzero-skills 创建一个用户管理 api..."
```

**觉得有用？请给 jzero 一个 ⭐ Star，支持我们继续改进！**

GitHub: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)
