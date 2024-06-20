---
title: 模版
icon: puzzle-piece
star: true
order: 5.2
category: 开发
tag:
  - Guide
---

jzero gen 生成代码, 仍然是依赖于 goctl 工具.

```shell
# 初始化的模板位置: ~/.jzero/$Version 下, 可以修改模板后再进行新建项目
jzero template init
# 或者 jzero template init --home .template

# 如果仍需要扩展 go-zero 的 template
goctl template init --home .template/go-zero
```

## 将当前项目转化为 jzero 的模板

```shell
jzero template build --name my_template

jzero new my_template_project --branch my_template
```