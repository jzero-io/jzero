---
title: Custom CLI plugins
icon: /icons/arcticons-game-plugins.svg
star: true
order: 0.15
---

If built-in commands are not enough, you can extend `jzero` with external executables instead of modifying the main binary.

This is useful for team-specific scaffolding, internal release workflows, deployment helpers, or any command that should feel like a native `jzero` subcommand.

## Discovery rules

When `jzero` receives an unknown command, it searches `PATH` for matching plugin executables:

* `jzero hello` -> `jzero-hello`
* `jzero foo bar` -> first tries `jzero-foo-bar`, then falls back to `jzero-foo`
* After a plugin is matched, the remaining arguments are passed through to the plugin unchanged
* The current environment variables are also forwarded to the plugin process

A plugin only needs two requirements:

* The file name starts with `jzero-`
* The file is executable and available in `PATH`

:::tip
Put plugin-specific flags after the plugin command, for example `jzero hello --name codex`.
:::

## Minimal example

The plugin can be written in Go, shell, or any language that can produce an executable in `PATH`.

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

## Read `desc` metadata in Go plugins

If your plugin is implemented in Go, you can also reuse `github.com/jzero-io/jzero/cmd/jzero/pkg/plugin`.

This does not replace external plugin discovery. `jzero` still discovers your binary through the `jzero-*` naming rule. The extra package is for reading parsed project metadata inside the plugin process.

`plugin.New()` scans the current working directory and attempts to parse:

* `desc/api`
* `desc/proto`
* `desc/sql`

It returns a `Metadata` value whose `Desc` field contains:

* `Desc.Api.SpecMap`: parsed API specs keyed by source file path
* `Desc.Proto.SpecMap`: parsed Proto specs keyed by source file path
* `Desc.Model.SpecMap`: parsed SQL table specs keyed by table name

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
`plugin.New()` reads from the plugin process's current working directory, so it is typically used when your plugin is executed inside a jzero project root.
:::

## Multi-level commands

You can map multiple command levels to a single executable name.

```bash
# jzero foo bar baz
# jzero will try jzero-foo-bar first
# if not found, it falls back to jzero-foo
# this usually means subcommands like "bar baz" are handled by jzero-foo itself
```

This allows you to organize team commands in a natural way, such as `jzero release publish` or `jzero company bootstrap`.

## Naming notes

Inside each command segment, `jzero` normalizes `-` to `_` before lookup.

For example:

* `jzero my-cmd` -> executable name `jzero-my_cmd`

To keep naming predictable, prefer simple command names or use `_` in the plugin executable when your command segment contains `-`.

## Recommended workflow

1. Build or place the plugin executable in a directory that is already in `PATH`
2. Follow the `jzero-<command>` naming rule
3. Add help output in the plugin itself, then use `jzero <command> --help` to view usage

Plugins are discovered dynamically, so they are not part of the built-in static command list printed by `jzero --help`.
