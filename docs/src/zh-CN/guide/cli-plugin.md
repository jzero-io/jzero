---
title: 自定义 jzero CLI 插件
icon: /icons/arcticons-game-plugins.svg
star: true
order: 0.15
---

如果内置命令不够用，可以通过外部可执行文件扩展 `jzero`，而不需要修改主二进制。

这适合承载团队内部脚手架、发布流程、部署辅助命令，或者任何希望表现成原生命令的自定义能力。

## 发现规则

当 `jzero` 收到一个未知命令时，会到 `PATH` 中查找匹配的插件可执行文件：

* `jzero hello` -> `jzero-hello`
* `jzero foo bar` -> 优先尝试 `jzero-foo-bar`，找不到时再回退到 `jzero-foo`
* 命中插件后，剩余参数会原样透传给插件
* 当前环境变量也会一并传递给插件进程

一个插件只需要满足两个条件：

* 文件名以 `jzero-` 开头
* 文件本身可执行，并且位于 `PATH` 中

:::tip
插件自己的参数请放在插件命令之后，例如 `jzero hello --name codex`。
:::

## 最小示例

插件可以用 Go、Shell，或者任何能够生成 `PATH` 内可执行文件的语言来实现。

```bash
mkdir -p ~/.local/bin

cat > ~/.local/bin/jzero-hello <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

name="${1:-world}"
printf 'hello, %s\n' "$name"
EOF

chmod +x ~/.local/bin/jzero-hello
export PATH="$HOME/.local/bin:$PATH"

jzero hello codex
# hello, codex
```

## 在 Go 插件中读取 `desc` 元数据

如果你的插件是用 Go 实现的，还可以复用 `github.com/jzero-io/jzero/cmd/jzero/pkg/plugin`。

这并不会替代前面的外部插件发现机制。`jzero` 仍然通过 `jzero-*` 命名规则发现并执行你的插件二进制；这个包解决的是插件进程内部如何读取并解析项目描述文件。

`plugin.New()` 会基于当前工作目录，尝试解析：

* `desc/api`
* `desc/proto`
* `desc/sql`

返回的 `Metadata` 中，`Desc` 字段包含：

* `Desc.Api.SpecMap`：按源文件路径组织的 API 解析结果
* `Desc.Proto.SpecMap`：按源文件路径组织的 Proto 解析结果
* `Desc.Model.SpecMap`：按表名组织的 SQL 解析结果

```go
package main

import (
	"fmt"

	jplugin "github.com/jzero-io/jzero/cmd/jzero/pkg/plugin"
)

func main() {
	metadata, err := jplugin.New()
	if err != nil {
		panic(err)
	}

	fmt.Printf("api files: %d\n", len(metadata.Desc.Api.SpecMap))
	fmt.Printf("proto files: %d\n", len(metadata.Desc.Proto.SpecMap))
	fmt.Printf("sql tables: %d\n", len(metadata.Desc.Model.SpecMap))
}
```

:::tip
`plugin.New()` 读取的是插件进程当前工作目录下的内容，因此通常应在 jzero 项目根目录中执行插件。
:::

## 多级命令

你可以把多级命令路径映射成一个插件可执行文件名。

```bash
# jzero foo bar baz
# jzero 会优先尝试 jzero-foo-bar
# 如果没找到，则会回退到 jzero-foo
# 这通常意味着后续的 "bar baz" 子命令由 jzero-foo 自己继续处理
```

这样可以自然地组织团队命令，例如 `jzero release publish` 或 `jzero company bootstrap`。

## 命名注意事项

在每个命令段内部，`jzero` 会在查找前把 `-` 归一化成 `_`。

例如：

* `jzero my-cmd` -> 对应的可执行文件名应为 `jzero-my_cmd`

为了减少歧义，建议优先使用简单命令名；如果命令段中必须包含 `-`，则在插件文件名中使用 `_`。

## 推荐工作流

1. 将插件可执行文件构建或放置到已经在 `PATH` 中的目录
2. 按照 `jzero-<command>` 规则命名
3. 在插件内部实现自己的帮助输出，然后通过 `jzero <command> --help` 查看使用方式

插件是动态发现的，因此不会出现在 `jzero --help` 打印的内置静态命令列表中。
