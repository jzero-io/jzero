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

## 实战: 构建属于自己的模版

:::tip 可以将当前任意项目转换成 jzero 模板, 这非常 cool!
:::

```bash
# 新增一个 api 项目
jzero new simpleapi
# 进入项目
cd simpleapi
# 新增一个 api
jzero add api helloworld
# 生成代码
jzero gen

# 将当前项目构建为模版, 并保存到 $HOME/.jzero/templates/local/myapi 下
jzero template build --name myapi

# 此时就可以使用你自己构建的模板了, 你会发现生成的项目自动拥有了 helloworld api 了.
jzero new mysimpleapi --local myapi

# 但是你发现该模板仅允许本地使用, 为了达到通用的效果
# 你可以在远程仓库如 github 创建一个 templates 仓库(假设为 https://github.com/jzero-io/templates)
# 然后将 $HOME/.jzero/templates/local/myapi 下的内容放到仓库中, 并上传到 myapi 分支
jzero new project_name --remote https://github.com/jzero-io/templates --branch myapi
```

模板结构如下:

```bash
$ tree ~/.jzero/templates/local/myapi
└── app
    ├── Dockerfile.tpl
    ├── README.md.tpl
    ├── cmd
    │   ├── root.go.tpl
    │   ├── server.go.tpl
    │   └── version.go.tpl
    ├── desc
    │   ├── api
    │   │   ├── helloworld.api.tpl
    │   │   └── version.api.tpl
    │   └── swagger
    │       ├── helloworld.swagger.json.tpl
    │       ├── swagger.json.tpl
    │       └── version.swagger.json.tpl
    ├── etc
    │   └── etc.yaml.tpl
    ├── go.mod.tpl
    ├── internal
    │   ├── config
    │   │   └── config.go.tpl
    │   ├── custom
    │   │   └── custom.go.tpl
    │   ├── handler
    │   │   ├── helloworld
    │   │   │   └── helloworld_compact.go.tpl
    │   │   ├── routes.go.tpl
    │   │   └── version
    │   │       └── version.go.tpl
    │   ├── logic
    │   │   ├── helloworld
    │   │   │   └── create.go.tpl
    │   │   └── version
    │   │       └── version.go.tpl
    │   ├── middleware
    │   │   ├── middleware.go.tpl
    │   │   ├── response.go.tpl
    │   │   └── validator.go.tpl
    │   ├── svc
    │   │   ├── config.go.tpl
    │   │   ├── middleware.go.tpl
    │   │   └── servicecontext.go.tpl
    │   └── types
    │       ├── helloworld
    │       │   └── types.go.tpl
    │       ├── types.go.tpl
    │       └── version
    │           └── types.go.tpl
    ├── main.go.tpl
    └── plugins
        └── plugins.go.tpl
```