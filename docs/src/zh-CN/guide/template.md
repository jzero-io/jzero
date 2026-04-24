---
title: 模版指南
icon: /icons/vscode-icons-folder-type-template.svg
star: true
order: 5.3
---

## 模版初始化

将jzero内嵌模版或者远程仓库的模版初始化到本地磁盘。

```shell
# 初始化jzero内嵌模板到 $HOME/.jzero/templates/$Version 下, 可以修改模板后再进行新建项目
jzero template init
# 或者初始化模板到当前项目的 .template, jzero gen 默认会优先读取当前项目的 .template 作为模板 home
jzero template init --output .template
# 初始化远程仓库模板到 $HOME/.jzero/templates/remote 下, 如 gateway, 
jzero template init --branch gateway

# 如果仍需要扩展 go-zero 的 template
goctl template init --home .template/go-zero
```

## 使用自定义模版初始化项目

* 指定远程仓库模板

```shell
jzero new project_name --remote repo_to_your_templates --branch template_branch
# 从缓存获取远程模板
jzero new project_name --remote repo_to_your_templates --branch template_branch --cache
```

* 使用本地模版

```shell
jzero new project_name --local template_name
```

* 使用路径模版

```shell
jzero new project_name --home path_to_template
```

## 模版渲染与变量

`jzero new` 在生成项目时，会同时渲染模板内容和模板路径：

* `.tpl` 文件内容会作为 Go `text/template` 渲染
* 文件名和目录名也会作为模板渲染，因此路径里同样可以使用 `{{ .APP }}`、`{{ .Module }}` 这类变量
* 如果文件名以 `.tpl.tpl` 结尾，只会去掉一层 `.tpl`，文件内容不会再次渲染，适合需要保留模板原文的场景

例如下面这个模板文件路径：

```text
internal/{{ .APP | lower }}/{{ FormatStyle .Style "service_context.go.tpl" }}
```

在项目创建时会被渲染成真实目录和文件名。

### 内置变量

执行 `jzero new` 时，jzero 会向模板注入以下内置变量：

| 变量 | 类型 | 说明 |
| --- | --- | --- |
| `APP` | `string` | 项目名，来自 `jzero new <name>` 或 `--name` |
| `Module` | `string` | Go module 名称，来自 `--module`，未指定时默认与项目名一致 |
| `GoVersion` | `string` | 当前 Go 版本 |
| `GoArch` | `string` | 当前架构，如 `amd64`、`arm64` |
| `DirName` | `string` | 输出目录名 |
| `Style` | `string` | 文件命名风格，默认 `gozero` |
| `Features` | `[]string` | `jzero new --features` 传入的特性列表 |
| `Serverless` | `bool` | 是否以 serverless 模式创建项目 |

例如：

```text
module {{ .Module }}

{{ if has "model" .Features }}
// enable model feature
{{ end }}

{{ if .Serverless }}
// serverless mode
{{ end }}
```

:::tip
`jzero template build` 会自动把项目中的 `go.mod` module，以及 Go 代码里引用当前项目的 import 路径，改写成 `{{ .Module }}`。因此通过 `jzero template build` 构建出来的模板，至少可以直接复用 `Module` 变量。
:::

### 内置函数

模板底层使用 Go `text/template`，除了 `and`、`or`、`not`、`index` 这类内置函数，还可以直接使用 [sprig](https://masterminds.github.io/sprig/) 提供的很多常用函数，例如 `lower`、`upper`、`default`、`has`、`dict` 等。除此之外，jzero 还额外注册了以下函数：

| 函数 | 说明 |
| --- | --- |
| `FirstUpper(s)` | 首字母转大写 |
| `FirstLower(s)` | 首字母转小写 |
| `ToCamel(s)` | 将 `foo-bar`、`foo_bar`、`foo/bar` 转成驼峰 |
| `FormatStyle(style, name)` | 按 `--style` 对文件名进行风格转换 |
| `VersionCompare(v1, op, v2)` | 版本比较，支持 `>、<、>=、<=` |

例如：

```text
{{ .APP | ToCamel | FirstUpper }}
{{ FormatStyle .Style "service_context.go.tpl" }}
{{ if (VersionCompare .GoVersion ">=" "1.24") }}toolchain go1.24.0{{ end }}
```

### 注入自定义模板变量

可以通过全局参数 `--register-tpl-val key=value` 注入额外模板变量。注入后的值会合并到当前模板数据中，因此既可以在模板内容中使用，也可以在模板路径中使用。

```shell
jzero new myapi --local myapi \
  --register-tpl-val company=acme \
  --register-tpl-val owner=platform
```

模板中可以直接访问：

```text
# {{ .APP }}
Company: {{ .company }}
Owner: {{ .owner }}
```

也可以用于路径：

```text
internal/{{ .company }}/banner.txt.tpl
```

如果你希望长期复用这些变量，也可以写到 `.jzero.yaml`：

```yaml
register-tpl-val:
  - company=acme
  - owner=platform
```

注意：

* 自定义注入变量与内置变量同名时，会覆盖原来的值
* 当前参数按 `key=value` 解析，建议 value 中不要再包含 `=`
* `--register-tpl-val` 是全局参数，不只 `jzero new` 可用；其他使用 jzero 模板渲染的命令也会额外合并这些变量，但不同命令自带的内置变量并不完全相同

## 实战: 构建属于自己的模版

:::tip 可以将当前任意项目转换成 jzero 模板, 这非常 cool!
:::

```bash
# 新增一个 api 项目
jzero new simpleapi
# 进入项目
cd simpleapi
# 新增一个 api
jzero add api helloworld
# 生成代码
jzero gen

# 将当前项目构建为模版, 并保存到 $HOME/.jzero/templates/local/myapi 下
jzero template build --name myapi

# 此时就可以使用你自己构建的模板了, 你会发现生成的项目自动拥有了 helloworld api 了.
jzero new mysimpleapi --local myapi

# 但是你发现该模板仅允许本地使用, 为了达到通用的效果
# 你可以在远程仓库如 github 创建一个 templates 仓库(假设为 https://github.com/jzero-io/templates)
# 然后将 $HOME/.jzero/templates/local/myapi 下的内容放到仓库中, 并上传到 myapi 分支
jzero new project_name --remote https://github.com/jzero-io/templates --branch myapi
```

模板结构如下:

```bash
$ tree ~/.jzero/templates/local/myapi
└── app
    ├── Dockerfile.tpl
    ├── README.md.tpl
    ├── cmd
    │   ├── root.go.tpl
    │   ├── server.go.tpl
    │   └── version.go.tpl
    ├── desc
    │   ├── api
    │   │   ├── helloworld.api.tpl
    │   │   └── version.api.tpl
    │   └── swagger
    │       ├── helloworld.swagger.json.tpl
    │       ├── swagger.json.tpl
    │       └── version.swagger.json.tpl
    ├── etc
    │   └── etc.yaml.tpl
    ├── go.mod.tpl
    ├── internal
    │   ├── config
    │   │   └── config.go.tpl
    │   ├── custom
    │   │   └── custom.go.tpl
    │   ├── handler
    │   │   ├── helloworld
    │   │   │   └── helloworld_compact.go.tpl
    │   │   ├── routes.go.tpl
    │   │   └── version
    │   │       └── version.go.tpl
    │   ├── logic
    │   │   ├── helloworld
    │   │   │   └── create.go.tpl
    │   │   └── version
    │   │       └── version.go.tpl
    │   ├── middleware
    │   │   ├── middleware.go.tpl
    │   │   ├── response.go.tpl
    │   │   └── validator.go.tpl
    │   ├── svc
    │   │   ├── config.go.tpl
    │   │   ├── middleware.go.tpl
    │   │   └── servicecontext.go.tpl
    │   └── types
    │       ├── helloworld
    │       │   └── types.go.tpl
    │       ├── types.go.tpl
    │       └── version
    │           └── types.go.tpl
    ├── main.go.tpl
    └── plugins
        └── plugins.go.tpl
```
