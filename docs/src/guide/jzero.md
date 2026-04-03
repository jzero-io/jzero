---
title: Mastering jzero
icon: /icons/catppuccin-astro-config.svg
star: true
order: 0.1
---

## About Configuration

* Supports controlling various parameters through configuration file .jzero.yaml
* Supports controlling various parameters through flag
* Supports controlling various parameters through environment variables
* Supports controlling various parameters through combination of above methods, priority from high to low: environment variables > flag > configuration file

Example: `jzero gen --style go_zero` corresponds to `.jzero.yaml` content

::: code-tabs#yaml
@tab .jzero.yaml
```yaml
gen:
  git-change: true
```
:::

`jzero gen` + `.jzero.yaml` = `jzero gen --git-change=true`

For environment variable usage, need to add prefix `JZERO_`, such as `JZERO_GEN_GIT_CHANGE`

`JZERO_GEN_GIT_CHANGE=go_zero jzero gen` = `jzero gen --git-change=true`

Environment variable definition supports using configuration file, default is `.jzero.env.yaml`

Example:

::: code-tabs#yaml
@tab .jzero.env.yaml
```yaml
JZERO_GEN_GIT_CHANGE: true
```
:::

### Subcommands

For subcommand configuration, such as: `jzero gen zrpcclient --output client` corresponds to `.jzero.yaml` content

::: code-tabs#yaml
@tab .jzero.yaml
```yaml
gen:
  zrpcclient:
    output: client
```
:::

`jzero gen zrpcclient` + `.jzero.yaml` = `jzero gen zrpcclient --output client`

Also supports environment variable configuration `JZERO_GEN_ZRPCCLIENT_NAME`

::: code-tabs#yaml
@tab .jzero.env.yaml
```yaml
JZERO_GEN_ZRPCCLIENT_OUTPUT: client
```
:::

`jzero gen zrpcclient` + `.jzero.env.yaml` = `jzero gen zrpcclient --output client`

## Set working directory

```shell
jzero gen -w /path/to
```

## Set quiet mode

```shell
jzero gen --quiet
```

## Set debug mode

```shell
jzero gen --debug
```
