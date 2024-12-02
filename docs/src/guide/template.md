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
jzero template build --name template_name
```

生成的模版将默认放在 `$HOME/.jzero/templates` 下。

## 使用模版构建项目

1. 使用远程模版

- 使用默认仓库：`https://github.com/jzero-io/templates`

```shell
jzero new project_name --branch template_branch
```

- 指定远程仓库

```shell
jzero new project_name --remote repo_to_your_templates --branch template_branch 
```

2. 使用本地缓存的远程模版

```shell
jzero new project_name --branch --cache template_branch
```
本地缓存的模版也在 `$HOME/.jzero/templates` 下。

3. 使用本地模版

```shell
jzero new project_name --local template_name
```

4. 使用特定模版

```shell
jzero new project_name --home path_to_template
```
