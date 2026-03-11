---
title: Generate server code
icon: vscode-icons:folder-type-api-opened
order: 4
---

jzero code generation command is extremely simple, only need `jzero gen` to automatically recognize all descriptor files/configurations and complete code generation.

After adding descriptor files with the `jzero add` command from the previous document, execute `jzero gen` to see the generated files.

## Generate code

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen
```

@tab Docker

```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
```
:::

## Generate code based on git changes

::: tip Get new/modified descriptor files based on git status -su
:::

```shell
jzero gen --git-change
```

## Generate code for specific desc

```shell
jzero gen --desc desc/api/xx.api
jzero gen --desc desc/proto/xx.proto
jzero gen --desc desc/sql/xx.sql
```

## Ignore specific desc when generating code

```shell
jzero gen --desc-ignore desc/api/xx.api
jzero gen --desc-ignore desc/proto/xx.proto
jzero gen --desc-ignore desc/sql/xx.sql
```

For more usage, see: [jzero guide](../guide/jzero.md)
