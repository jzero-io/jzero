---
title: 模版
icon: vscode-icons:folder-type-template
star: true
order: 5.3
category: 开发
tag:
  - Guide
---

jzero gen 生成代码, 仍然是依赖于 goctl 工具.

```shell
# 初始化的模板位置: ~/.jzero/$Version 下, 可以修改模板后再进行新建项目
jzero template init
# 或者初始化模板到当前项目, jzero gen 默认会优先读取当前项目的 .template 作为模板 home
jzero template init --output .template
# 初始化远程仓库模板. 如 gateway
jzero template init --output .template --branch gateway

# 如果仍需要扩展 go-zero 的 template
goctl template init --home .template/go-zero
```

## 将当前项目转化为 jzero 的模板

```shell
jzero template build --name my_template

jzero new my_template_project --branch my_template
```