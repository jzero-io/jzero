---
title: 插件指南
icon: arcticons:game-plugins
star: true
order: 5.4
---

jzero 支持插件化机制, 可以方便的进行插件的安装和卸载操作. 

核心点在于**多模块协同开发**, 最终编译成**单体服务部署**.

## 新增插件(以 api 项目为例)

```bash
# 新增 api 项目
jzero new simpleapi
# 进入项目目录
cd simpleapi
# 新增 api 项目插件(独立 go module)
jzero new your_plugin --frame api --serverless
# 新增 api 项目插件(与主服务 simpleapi 共用 go module )
jzero new your_mono_plugin --frame api --serverless --mono
# 执行 serverless build, 主服务接管插件路由(plugins/plugins.go)
jzero serverless build
# 下载依赖
go mod tidy
# 大单体编译产物
go build
```

## 卸载插件

```shell
# 卸载所有, 主服务不再接管插件路由
jzero serverless delete

# 卸载指定插件
jzero serverless delete --plugin <plugin-name>

# 重新编译
go build
```

## 项目结构

```bash
simpleapi
├── Dockerfile
├── README.md
├── cmd
│   ├── root.go
│   ├── server.go
│   └── version.go
├── desc
│   ├── api
│   │   └── version.api
│   └── swagger
│       ├── swagger.json
│       └── version.swagger.json
├── etc
│   └── etc.yaml
├── go.mod
├── go.sum
├── go.work
├── go.work.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── custom
│   │   └── custom.go
│   ├── handler
│   │   ├── routes.go
│   │   └── version
│   │       └── version.go
│   ├── logic
│   │   └── version
│   │       └── version.go
│   ├── middleware
│   │   ├── middleware.go
│   │   ├── response.go
│   │   └── validator.go
│   ├── svc
│   │   ├── config.go
│   │   ├── middleware.go
│   │   └── servicecontext.go
│   └── types
│       ├── types.go
│       └── version
│           └── types.go
├── main.go
└── plugins
    ├── plugins.go
    ├── your_mono_plugin
    │   ├── Dockerfile
    │   ├── README.md
    │   ├── cmd
    │   │   ├── root.go
    │   │   ├── server.go
    │   │   └── version.go
    │   ├── etc
    │   │   └── etc.yaml
    │   ├── internal
    │   │   ├── config
    │   │   │   └── config.go
    │   │   ├── custom
    │   │   │   └── custom.go
    │   │   ├── handler
    │   │   │   └── routes.go
    │   │   ├── middleware
    │   │   │   ├── middleware.go
    │   │   │   ├── response.go
    │   │   │   └── validator.go
    │   │   └── svc
    │   │       ├── config.go
    │   │       ├── middleware.go
    │   │       └── servicecontext.go
    │   ├── main.go
    │   └── serverless
    │       └── serverless.go
    └── your_plugin
        ├── Dockerfile
        ├── README.md
        ├── cmd
        │   ├── root.go
        │   ├── server.go
        │   └── version.go
        ├── etc
        │   └── etc.yaml
        ├── go.mod
        ├── internal
        │   ├── config
        │   │   └── config.go
        │   ├── custom
        │   │   └── custom.go
        │   ├── handler
        │   │   └── routes.go
        │   ├── middleware
        │   │   ├── middleware.go
        │   │   ├── response.go
        │   │   └── validator.go
        │   └── svc
        │       ├── config.go
        │       ├── middleware.go
        │       └── servicecontext.go
        ├── main.go
        └── serverless
            └── serverless.go
```