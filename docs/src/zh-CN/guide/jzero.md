---
title: 玩转 jzero
icon: catppuccin:astro-config
star: true
order: 0.1
---

## 关于配置

* 支持通过配置文件 .jzero.yaml 控制各种参数
* 支持通过 flag 控制各种参数
* 支持通过环境变量控制各种参数
* 支持通过以上组合的方式控制各种参数, 优先级从高到低为: 环境变量  > flag  > 配置文件

如: `jzero gen --style go_zero` 对应 `.jzero.yaml` 内容

::: code-tabs#yaml
@tab .jzero.yaml
```yaml
gen:
  git-change: true
```
:::

即 `jzero gen` + `.jzero.yaml` = `jzero gen --git-change=true`

对于环境变量的使用, 需要增加前缀 `JZERO_`, 如 `JZERO_GEN_GIT_CHANGE`

即 `JZERO_GEN_GIT_CHANGE=go_zero jzero gen` = `jzero gen --git-change=true`

环境变量的定义支持使用配置文件, 默认为 `.jzero.env.yaml`

如:

::: code-tabs#yaml
@tab .jzero.env.yaml
```yaml
JZERO_GEN_GIT_CHANGE: true
```
:::

### 子命令

对于子命令的配置, 如: `jzero gen zrpcclient --output client` 对应 `.jzero.yaml` 内容

::: code-tabs#yaml
@tab .jzero.yaml
```yaml
gen:
  zrpcclient:
    output: client
```
:::

`jzero gen zrpcclient` + `.jzero.yaml` = `jzero gen zrpcclient --output client`

同样支持环境变量的配置 `JZERO_GEN_ZRPCCLIENT_NAME`

::: code-tabs#yaml
@tab .jzero.env.yaml
```yaml
JZERO_GEN_ZRPCCLIENT_OUTPUT: client
```
:::

`jzero gen zrpcclient` + `.jzero.env.yaml` = `jzero gen zrpcclient --output client`

## 设置工作目录

```shell
jzero gen -w /path/to
```

## 设置 quiet 模式

```shell
jzero gen --quiet
```

## 设置 debug 模式

```shell
jzero gen --debug
```