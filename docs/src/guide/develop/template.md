---
title: 模版定制化
icon: puzzle-piece
star: true
order: 2
category: 开发
tag:
  - Guide
---

jzero gen 生成代码, 仍然是依赖于 goctl 工具.

:::tip go-zero 的模版必须放在项目根路径 .template/go-zero 位置, 否则不会生效. 

另外由于做了定制化, 部分 template 修改将不会生效.
:::

```shell
goctl template init --home .template/go-zero
```