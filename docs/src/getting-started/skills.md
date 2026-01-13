---
title: Skills
icon: logos:claude-icon
order: 6
---

# jzero Skills

`jzero skills` 命令用于将 jzero 框架的 AI 技能模板复制到项目目录中，使 AI 助手（如 Claude）能够更好地理解和编写 jzero 代码。

## 基础使用

### 命令语法

```bash
jzero skills [flags]
```

### 参数说明

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--output` | `-o` | `.claude/skills` | 输出目录 |

### 示例

```bash
# 使用默认输出目录 .claude/skills
jzero skills

# 指定自定义输出目录
jzero skills -o /path/to/custom/skills
```

## 使用流程

### 1. 执行命令

在项目根目录执行：

```bash
jzero skills
```

执行成功后会显示：

```
Skills templates copied successfully to: /path/to/project/.claude/skills
```

### 2. 查看生成的文件结构

```
.claude/skills/
└── jzero-skills/
    ├── SKILL.md                    # 技能入口文件
    └── references/                 # 详细模式文档
        ├── best-practices.md       # 数据库最佳实践
        ├── condition-builder.md    # 条件构建器使用
        ├── crud-operations.md      # CRUD 操作示例
        ├── database-connection.md  # 数据库连接配置
        ├── database-patterns.md    # 数据库操作模式
        ├── model-generation.md     # Model 生成指南
        └── rest-api-patterns.md    # REST API 模式
```

### 3. AI 助手自动加载

当你使用 Claude Code 或其他支持 Skills 的 AI 工具时，这些技能模板会被自动识别和加载，AI 助手将能够：

- 理解 jzero 的三层架构（Handler → Logic → Model）
- 正确使用 `condition.NewChain()` API 构建查询条件
- 遵循 jzero 的代码生成流程
- 应用框架最佳实践和模式

## Skills 的作用

### 1. 提供 jzero 框架知识

AI 助手通过学习 Skills 文档，能够掌握：

- **REST API 开发**：Handler/Logic/Context 三层架构、请求验证、错误处理、中间件实现
- **数据库操作**：SQL 操作、Redis 缓存、连接池优化
- **jzero 特性**：基于 git 变更的代码生成、灵活配置、自定义模板

### 2. 确保代码规范

Skills 中定义了关键的开发规范：

```go
// ✅ 正确 - 使用 condition.NewChain() API
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/user"
)

conditions := condition.NewChain().
    Equal(user.Username, req.Username).
    Greater(user.Age, 18).
    Build()

// ❌ 错误 - 不要使用 condition.New()
conditions := condition.New(
    condition.Condition{Field: user.Username, Operator: condition.Equal, Value: req.Username},
)
```

### 3. 指导开发流程

Skills 明确定义了代码生成的工作流：

```bash
# 1. 修改 .api/.proto/.sql 文件
# 2. 运行对应的生成命令
jzero gen --desc desc/api/your_file.api
jzero gen --desc desc/proto/your_file.proto
jzero gen --desc desc/sql/your_file.sql

# 3. 在生成的文件中实现业务逻辑
# 4. 测试并提交
```

### 4. 提供模式参考

Skills 包含详细的开发模式文档，涵盖：

- CRUD 完整示例
- 条件查询构建
- 中间件实现
- 事务处理
- 分表分库策略
- 缓存策略

## 注意事项

1. **目录位置**：默认输出到 `.claude/skills`，这是 Claude Code 的标准技能目录
2. **覆盖更新**：如果目录已存在，执行命令会覆盖现有文件
3. **版本兼容**：Skills 针对 jzero 1.1+ 版本，保持与 go-zero 1.5+ 兼容
4. **AI 工具支持**：主要用于 Claude Code 等 AI 辅助开发工具
