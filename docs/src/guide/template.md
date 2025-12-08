---
title: 模版指南
icon: vscode-icons:folder-type-template
star: true
order: 5.3
---

## 模版初始化

将jzero内嵌模版或者远程仓库的模版初始化到本地磁盘。

```shell
# 初始化jzero内嵌模板到 $HOME/.jzero/templates/$Version 下, 可以修改模板后再进行新建项目
jzero template init
# 或者初始化模板到当前项目的 .template, jzero gen 默认会优先读取当前项目的 .template 作为模板 home
jzero template init --output .template
# 初始化远程仓库模板到 $HOME/.jzero/templates/remote 下, 如 gateway, 
jzero template init --branch gateway

# 如果仍需要扩展 go-zero 的 template
goctl template init --home .template/go-zero
```

## 构建属于自己的模版

```shell
# 将当前项目构建为模版，并保存到 $HOME/.jzero/templates/local 下
jzero template build --name template_name
```

## 使用自定义模版初始化项目

* 指定远程仓库模板

```shell
jzero new project_name --remote repo_to_your_templates --branch template_branch
# 从缓存获取远程模板
jzero new project_name --remote repo_to_your_templates --branch template_branch --cache
```

* 使用本地模版

```shell
jzero new project_name --local template_name
```

* 使用路径模版

```shell
jzero new project_name --home path_to_template
```