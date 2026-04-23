---
title: AI 时代的利器 CLI 工具指南
icon: /icons/streamline-ultimate-blog-blogger-logo.svg
---

在 AI 时代，命令行工具（CLI）正在经历一场**复古潮流**的复兴。为什么？因为 **CLI 工具对 AI Agent 天然友好**！

与图形界面相比，CLI 工具具有结构化的输入输出、清晰的帮助文档、可预测的行为模式——这些特性使 AI Agent 能够更容易地理解、学习和自动化使用 CLI 工具。

但你是否曾为编写一个功能完善的 CLI 工具而烦恼：

- 命令参数解析繁琐？
- 配置管理混乱？
- 帮助文档难以维护？
- 想要插件扩展却无从下手？

---

💡 **在此之前，先介绍一下 jzero 模板市场**

jzero 提供了丰富的**官方模板**和**第三方模板**，帮助你快速构建各种类型的项目：

**🚀 内置模板**：

- **RPC 模板**：基于 Protocol Buffers 的高性能 gRPC 微服务框架
- **API 模板**：基于 API 描述语言的轻量级 RESTful API 服务框架
- **Gateway 模板**：高性能 API 网关，同时支持 gRPC 和 HTTP 协议

**📦 官方外置模板**：
- **CLI 模板**：具有常见 CLI 模式的命令行应用程序模板（今天的主角！）
- **API 模板**：专为 Vercel 部署优化的 API 模板
- **文档模板**：使用 VuePress Hope 主题的文档站点模板

**🌍 第三方模板**：
- 欢迎贡献你自己的模板，帮助更多开发者快速启动项目！

访问 **[jzero 模板市场](https://templates.jzero.io/)** 了解更多模板信息和使用指南。

![](https://oss.jaronnie.com/image-20260409190255415.png)

---

今天，我们将介绍如何使用 **jzero CLI 模板**快速构建专业的命令行工具！

![](https://oss.jaronnie.com/image-20260409190335125.png)

---

## 为什么选择 jzero CLI 模板？

jzero CLI 模板基于业界成熟的 **Cobra 框架**，提供了开箱即用的项目结构和最佳实践配置。相比从零开始搭建，使用 jzero CLI 模板能够：

✅ **快速启动**：一键生成完整项目结构，无需繁琐配置  
✅ **规范统一**：遵循行业标准，命令结构清晰易懂  
✅ **功能完备**：内置配置管理、插件系统、调试模式等企业级特性  
✅ **易于扩展**：插件化架构，轻松添加新功能  
✅ **AI 友好**：与 Claude、GPT 等 AI 工具完美配合，提升开发效率  

---

## 快速开始：1 分钟创建你的第一个 CLI 工具

```bash
# 1. 安装 jzero（如果尚未安装）
go install github.com/jzero-io/jzero/cmd/jzero@latest

# 2. 创建新的 CLI 项目
jzero new mycli --branch cli

cd mycli

# 3. 安装依赖并构建
go mod tidy
go build

# 4. 测试运行
./mycli version
```

输出示例：
```
mycli version 1.0.0 darwin/amd64
Go version go1.21.0
Git commit abc123
Build date 2024-01-01 12:00:00 +0000 UTC
```

就这么简单！你已经拥有了一个功能完整的 CLI 工具框架。

项目结构如下：

```
mycli/
├── main.go                    # 入口文件
├── internal/
│   ├── command/              # 命令实现
│   │   └── version/          # 版本命令
│   │       └── version.go
│   └── config/               # 配置管理
│       └── config.go
├── go.mod
└── go.sum
```

---

## 核心概念：命令系统的三层结构

jzero CLI 模板基于 Cobra 框架，采用清晰的**三层命令结构**：

```
Root Command（根命令）
├── Command（命令）
│   └── Sub Command（子命令）
```

### 1. Root Command（根命令）

**根命令**是 CLI 工具的入口点，定义了工具的基本信息、全局配置和顶层命令。

```go
// main.go
var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "我的 CLI 工具",
    Long:  `一个功能强大的命令行工具，帮助你提高工作效率`,
}
```

根命令的特性：
- ✅ 定义全局标志（如 `--debug`、`--config`）
- ✅ 注册顶层命令
- ✅ 提供工具的整体帮助信息

### 2. Command（命令）

**命令**是根命令下的直接子命令，代表主要的功能模块。

```bash
mycli version      # 版本命令 - 显示版本信息
mycli config       # 配置命令 - 管理配置
mycli plugin       # 插件命令 - 管理插件
mycli server       # 服务器命令 - 启动服务
```

命令的特点：
- ✅ 直接挂载在根命令下
- ✅ 可以有独立的标志和参数
- ✅ 可以包含子命令，形成命令树

### 3. Sub Command（子命令）

**子命令**是命令的下一级，用于实现更细分的功能。

```bash
# 配置命令的子命令
mycli config init      # 初始化配置
mycli config set       # 设置配置项
mycli config get       # 获取配置项
mycli config list      # 列出所有配置

# 插件命令的子命令
mycli plugin install   # 安装插件
mycli plugin remove    # 卸载插件
mycli plugin list      # 列出插件
mycli plugin update    # 更新插件
```

子命令的优势：
- ✅ 功能模块化，逻辑清晰
- ✅ 支持多层嵌套（如 `mycli config database connect`）
- ✅ 每个子命令可以独立开发和维护

### 命令示例对比

```bash
# Root Command
mycli                    # 执行根命令

# Command（一级命令）
mycli config             # 执行配置命令
mycli plugin             # 执行插件命令

# Sub Command（二级命令）
mycli config init        # 执行配置初始化子命令
mycli plugin install     # 执行插件安装子命令

# 更深层次的子命令（三级命令）
mycli server start       # 启动服务器
mycli server stop        # 停止服务器
mycli server status      # 查看服务器状态
```

### 标志（Flags）的类型

标志用于配置命令的行为，分为三种类型：

**局部标志**（Local Flags）- 仅对当前命令有效：
```go
Cmd.Flags().StringP("output", "o", "", "输出文件")
```

**持久化标志**（Persistent Flags）- 对当前命令及其所有子命令有效：
```go
rootCmd.PersistentFlags().BoolP("verbose", "v", false, "详细输出模式")
```

**全局标志**（Global Flags）- 对所有命令有效：
```bash
mycli --debug          # 启用调试模式
mycli --config file.yaml  # 指定配置文件
```

---

## 配置管理：灵活的多层次配置方案

jzero CLI 模板提供了强大的配置管理系统，支持三种配置方式的灵活组合：

### 配置优先级
```
命令行标志 > 环境变量 > 配置文件
```

### 1. 配置文件

默认配置文件位置：`~/.mycli.yaml`

```yaml
# 调试模式
debug: false

# 调试睡眠时间（秒）
debug-sleep-time: 0

# 问候配置
greet:
  name: 世界
```

### 2. 环境变量配置


jzero CLI 模板会自动将环境变量映射到配置字段，无需手动设置：

```bash
# 直接设置环境变量，自动映射到配置
export MYCLI_DEBUG=true
export MYCLI_DEBUG_SLEEP_TIME=5
export MYCLI_GREET_NAME="张三"
```

**环境变量命名规则**：

格式：`{APP_PREFIX}_{CONFIG_PATH}`

- `{APP_PREFIX}`：应用名前缀（大写），如 `MYCLI`、`JZERO`
- `{CONFIG_PATH}`：配置路径，`.` 和 `-` 替换为 `_`

映射示例：
| 配置字段 | 环境变量 |
|---------|----------|
| `config.C.Debug` | `MYCLI_DEBUG` |
| `config.C.DebugSleepTime` | `MYCLI_DEBUG_SLEEP_TIME` |
| `config.C.Greet.Name` | `MYCLI_GREET_NAME` |

### 3. 命令行标志

```bash
# 通过命令行覆盖配置
./mycli --debug
./mycli --config custom.yaml
```

---

## 统一配置管理

jzero CLI 模板提供统一的配置管理系统，所有配置通过 `internal/config/config.go` 管理，支持配置文件、环境变量和命令行标志三种方式，自动按优先级加载。

配置优先级：**命令行标志 > 环境变量 > 配置文件**

### 统一配置的优势

✅ **自动优先级处理**：Viper 自动处理配置优先级
✅ **类型安全**：使用 Go 结构体和 mapstructure 标签
✅ **环境变量支持**：自动将环境变量映射到配置字段
✅ **灵活扩展**：添加新配置字段非常简单
✅ **全局访问**：通过 `config.C` 全局变量在任何地方访问
✅ **命令特定配置**：支持为不同命令设置独立配置

---

## 添加自定义命令：四步完成

添加一个新命令需要四个步骤，按照实际开发流程依次进行：

### 第一步：创建命令文件

在 `internal/command/` 下创建新目录和命令文件：

```go
// internal/command/greet/greet.go
package greet

import (
    "fmt"
    "mycli/internal/config"
    "github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
    Use:   "greet",
    Short: "问候命令",
    Long:  `向用户发送友好的问候`,
    Run: func(cmd *cobra.Command, args []string) {
        // 从统一配置获取名字（映射到 greet.name）
        name := config.C.Greet.Name
        fmt.Printf("你好，%s！\n", name)
    },
}
```

### 第二步：定义命令标志

在命令的 `init()` 函数中定义 flags：

```go
func init() {
    // 定义短标志和长标志，设置默认值
    Cmd.Flags().StringP("name", "n", "世界", "指定问候对象的名字")
}
```

**重要说明**：
- 标志名称使用小写（如 `name`）
- 会自动映射到 `config.C.Greet.Name`
- 短标志 `-n` 是长标志 `--name` 的简写

### 第三步：添加配置字段

在 `internal/config/config.go` 中添加对应的配置结构：

```go
type Config struct {
    Debug bool `mapstructure:"debug"`
    DebugSleepTime int `mapstructure:"debug-sleep-time"`

    // 添加 greet 命令的配置
    Greet GreetConfig `mapstructure:"greet"`
}

// GreetConfig greet 命令的配置
type GreetConfig struct {
    Name string `mapstructure:"name"`
}
```

**配置映射规则**：
- 命令标志 `name` → 配置字段 `Greet.Name`
- 配置文件中的 `greet.name` → `config.C.Greet.Name`
- 环境变量 `MYCLI_GREET_NAME` → `config.C.Greet.Name`

### 第四步：注册命令

在 `main.go` 中导入并注册命令：

```go
import (
    "mycli/internal/command/greet"
    // 其他导入...
)

func init() {
    rootCmd.AddCommand(greet.Cmd)
}
```

### 测试使用

```bash
go build

# 方式1：使用默认值
./mycli greet
# 输出：你好，世界！

# 方式2：使用命令行标志
./mycli greet --name 张三
# 输出：你好，张三！

# 方式3：使用配置文件
echo "greet:" >> ~/.mycli.yaml
echo "  name: 李四" >> ~/.mycli.yaml
./mycli greet
# 输出：你好，李四！

# 方式4：使用环境变量
export MYCLI_GREET_NAME="王五"
./mycli greet
# 输出：你好，王五！
```

---

## 调试模式

jzero CLI 模板内置了完善的调试支持：

### 启用调试的三种方式

**方式一：配置文件**
```yaml
# ~/.mycli.yaml
debug: true
debug-sleep-time: 5  # 调试睡眠时间（秒）
```

**方式二：环境变量**
```bash
export MYCLI_DEBUG=true
./mycli
```

**方式三：命令行标志**
```bash
./mycli --debug
./mycli --debug --debug-sleep-time 5
```

### 调试功能特性

- **睡眠时间控制**：在关键步骤间暂停，方便 debug

---

## 插件系统

jzero CLI 模板支持强大的插件系统，让你的工具能够动态扩展功能。

### 插件命名规则

插件可执行文件以 `YOUR_APP-` 为前缀，例如：
- `mycli-git`
- `mycli-docker`
- `mycli-deploy`

### 插件自动发现

系统会自动从 PATH 中搜索以 `mycli-` 开头的可执行文件。

### 插件使用示例

```bash
# 安装插件到 PATH
sudo cp mycli-git /usr/local/bin/

# 直接使用插件
./mycli git status
./mycli docker build
```
## CLI 工具分发：使用 GoReleaser 和 GitHub Workflows

开发完成后，如何方便地分发你的 CLI 工具？**GoReleaser** 结合 **GitHub Actions** 可以实现自动化构建和发布！

### GoReleaser 配置

在项目根目录创建 `.goreleaser.yaml`：

```yaml
version: 2

before:
  hooks:
    - go mod tidy

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    dir: .
    id: mycli
    binary: mycli

archives:
  - formats: [tar.gz]
    # this name template makes the OS and Arch compatible with the results of `uname`.
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
      {{- if .Arm }}v{{ .Arm }}{{ end }}
    # use zip for windows archives
    format_overrides:
      - goos: windows
        formats: [zip]

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
```

### GitHub Actions 配置

创建 `.github/workflows/release.yml`：

> **注意**：使用前需要在 GitHub 仓库的 Secrets 中配置 `ACCESS_TOKEN`（Settings → Secrets and variables → Actions → New repository secret），需要有仓库写入权限。可以通过 Personal Access Tokens (PAT) 生成。
>
> ![](https://oss.jaronnie.com/image-20260410113251665.png)

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v4

      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Get version
        id: get_version
        run: echo ::set-output name=VERSION::${GITHUB_REF/refs\/tags\//}

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v7
        with:
          distribution: goreleaser
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.ACCESS_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v4

      - name: Log in to ghcr
        uses: docker/login-action@v4
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.ACCESS_TOKEN }}

      - name: Docker build and push
        uses: docker/build-push-action@v6
        with:
          registry: ghcr.io
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:latest
            ghcr.io/${{ github.repository }}:${{ steps.get_version.outputs.VERSION }}
```

### Dockerfile 配置

创建 `Dockerfile` 用于容器化部署：

```dockerfile
FROM alpine:latest

ENV CGO_ENABLED=0

LABEL \
    org.opencontainers.image.title="mycli" \
    org.opencontainers.image.description="My CLI tool" \
    org.opencontainers.image.url="https://github.com/yourusername/mycli" \
    org.opencontainers.image.documentation="https://github.com/yourusername/mycli#readme" \
    org.opencontainers.image.source="https://github.com/yourusername/mycli" \
    org.opencontainers.image.licenses="MIT" \
    maintainer="your-name <your-email@example.com>"

WORKDIR /app

COPY dist/mycli_linux_amd64_v1/mycli /dist/mycli_linux_amd64/mycli
COPY dist/mycli_linux_arm64_v8.0/mycli /dist/mycli_linux_arm64/mycli

RUN if [ $(go env GOARCH) = "amd64" ]; then \
      cp /dist/mycli_linux_amd64/mycli /usr/local/bin/mycli; \
    elif [ $(go env GOARCH) = "arm64" ]; then \
      cp /dist/mycli_linux_arm64/mycli /usr/local/bin/mycli; \
    fi

RUN apk update --no-cache \
    && apk add --no-cache tzdata ca-certificates \
    && rm -rf /dist

ENTRYPOINT ["mycli"]
```

### 发布流程

**1. 创建并推送标签**：

```bash
# 创建版本标签
git tag v1.0.0

# 推送标签到 GitHub
git push origin v1.0.0
```

**2. GitHub Actions 自动执行**：

- 触发 `.github/workflows/release.yml`
- GoReleaser 自动构建多平台二进制文件
- 生成 GitHub Release
- 上传构建产物和校验和
- 构建并推送 Docker 镜像到 GHCR

**3. 用户安装方式**：

**方式一：直接下载二进制文件**

```bash
# 从 GitHub Releases 下载
wget https://github.com/yourusername/mycli/releases/download/v1.0.0/mycli_linux_amd64.tar.gz

tar -xzf mycli_1.0.0_linux_amd64.tar.gz
sudo mv mycli /usr/local/bin/
```

**方式二：使用 go install 命令安装**

```bash
# 直接安装指定版本
go install github.com/yourusername/mycli@v1.0.0

# 安装最新版本
go install github.com/yourusername/mycli@latest
```

**方式三：使用 Docker 镜像**

```bash
# 直接运行
docker run --rm ghcr.io/yourusername/mycli:latest version

# 创建别名方便使用
alias mycli='docker run --rm -v $(pwd):/app -w /app ghcr.io/yourusername/mycli:latest'

# 然后就可以像本地安装一样使用
mycli version
mycli --help
```

或者创建一个 shell 脚本 `mycli-docker.sh`：

```bash
#!/bin/bash
docker run --rm -v "$(pwd)":/app -w /app ghcr.io/yourusername/mycli:latest "$@"
```

然后将其添加到 PATH：

```bash
chmod +x mycli-docker.sh
sudo mv mycli-docker.sh /usr/local/bin/mycli
```

### 自动化分发流程的优势

✅ **全自动化流程**：从代码标签到发布完成的全程自动化
✅ **多平台支持**：一次构建，支持 Linux、macOS、Windows 多架构
✅ **多格式分发**：二进制文件、Docker 镜像、go install 多种安装方式
✅ **版本管理**：自动注入版本信息，通过标签管理发布
✅ **安全性保障**：自动生成 SHA256 校验和，确保下载完整性
✅ **容器化部署**：自动构建并推送多架构 Docker 镜像到 GHCR

---

## 完整实例

为了帮助你更好地理解如何使用 jzero CLI 模板构建命令行工具，我们创建了一个完整的演示项目。

**项目地址**：[https://github.com/jaronnie/mycli](https://github.com/jaronnie/mycli)

### 项目特点

这个演示项目完全按照本文档的流程创建，包含了：

- ✅ **greet 命令**：演示如何添加自定义命令
- ✅ **统一配置管理**：展示配置文件、环境变量、命令行标志的使用
- ✅ **GoReleaser 配置**：完整的跨平台构建配置
- ✅ **GitHub Actions 工作流**：自动化发布流程
- ✅ **Dockerfile**：多架构容器化支持
- ✅ **完整文档**：详细的 README 和使用说明

### 快速体验

```bash
# 克隆项目
git clone https://github.com/jaronnie/mycli.git
cd mycli

# 安装依赖并构建
go mod tidy
go build

# 测试运行
./mycli version
./mycli greet
./mycli greet --name 张三
```

### 项目结构

```
mycli/
├── main.go                    # 入口文件
├── internal/
│   ├── command/              # 命令实现
│   │   ├── version/          # 版本命令
│   │   └── greet/            # 问候命令（自定义）
│   └── config/               # 配置管理
│       └── config.go
├── Dockerfile                 # Docker 配置
├── .goreleaser.yaml          # GoReleaser 配置
├── .github/workflows/
│   └── release.yml           # GitHub Actions
├── go.mod
├── go.sum
└── README.md                 # 项目文档
```

### 作为项目模板

你可以直接以 mycli 项目为基础，开发自己的 CLI 工具：

```bash
# Fork 或克隆项目
git clone https://github.com/jaronnie/mycli.git your-cli
cd your-cli

# 修改配置
# - 编辑 main.go 中的应用名称和描述
# - 修改 go.mod 中的模块路径
# - 添加你自己的命令
# - 更新 README.md

# 开始开发！
```

这个项目展示了 jzero CLI 模板的最佳实践，是学习 CLI 工具开发的绝佳起点。

![](https://oss.jaronnie.com/image-20260410115131704.png)

![](https://oss.jaronnie.com/image-20260410115221574.png)

---

## 相关资源

- **jzero GitHub**: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)
- **jzero 文档**: [https://docs.jzero.io](https://docs.jzero.io)
- **CLI 模板**: [https://templates.jzero.io/external/cli/](https://templates.jzero.io/external/cli/)
- **Cobra 文档**: [https://github.com/spf13/cobra](https://github.com/spf13/cobra)
- **Viper 文档**: [https://github.com/spf13/viper](https://github.com/spf13/viper)

---

**让 jzero CLI 模板成为你 AI 时代的得力助手！** 🚀

**觉得有用？请给 jzero 一个 ⭐ Star，支持我们继续改进！**
