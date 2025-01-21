---
title: jzero 概览
icon: grommet-icons:overview
order: 3
---

## 不同姿势使用 jzero

:::important 涨知识的小技巧
:::

* 支持通过配置文件 .jzero.yaml 控制各种参数(**强烈推荐在每个项目的根目录新建该文件**)
* 支持通过 flag 控制各种参数
* 支持通过环境变量控制各种参数
* 支持通过以上组合的方式控制各种参数, 优先级从高到低为 环境变量  > flag  > 配置文件

如: `jzero gen --style go_zero` 对应 .jzero.yaml 内容

```yaml
gen:
  style: go_zero
```

即 `jzero gen` + `.jzero.yaml` = `jzero gen --style go_zero`

对于环境变量的使用, 需要增加前缀 `JZERO_`, 如 `JZERO_GEN_STYLE`

即 `JZERO_GEN_STYLE=go_zero jzero gen` = `jzero gen --style go_zero`