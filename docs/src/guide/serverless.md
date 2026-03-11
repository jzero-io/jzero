---
title: Plugin guide
icon: arcticons:game-plugins
star: true
order: 5.4
---

jzero supports plugin mechanism, making it easy to install and uninstall plugins.

The key point is **multi-module collaborative development**, finally compiled into **monolithic service deployment**.

## Add plugin (using api project as example)

```bash
# Add new api project
jzero new simpleapi
# Enter project directory
cd simpleapi
# Add api project plugin (independent go module)
jzero new your_plugin --frame api --serverless
# Add api project plugin (share go module with main service simpleapi)
jzero new your_mono_plugin --frame api --serverless --mono
# Execute serverless build, main service takes over plugin routes (plugins/plugins.go)
jzero serverless build
# Download dependencies
go mod tidy
# Large monolithic build output
go build
```

## Uninstall plugin

```shell
# Uninstall all, main service no longer takes over plugin routes
jzero serverless delete

# Uninstall specific plugin
jzero serverless delete --plugin <plugin-name>

# Rebuild
go build
```

## Project structure

```bash
simpleapi
├── Dockerfile
├── README.md
├── cmd
│   ├── root.go
│   ├── server.go
│   └── version.go
├── desc
│   ├── api
│   │   └── version.api
│   └── swagger
│       ├── swagger.json
│       └── version.swagger.json
├── etc
│   └── etc.yaml
├── go.mod
├── go.sum
├── go.work
├── go.work.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── custom
│   │   └── custom.go
│   ├── handler
│   │   ├── routes.go
│   │   └── version
│   │       └── version.go
│   ├── logic
│   │   └── version
│   │       └── version.go
│   ├── middleware
│   │   ├── middleware.go
│   │   ├── response.go
│   │   └── validator.go
│   ├── svc
│   │   ├── config.go
│   │   ├── middleware.go
│   │   └── servicecontext.go
│   └── types
│       ├── types.go
│       └── version
│           └── types.go
├── main.go
└── plugins
    ├── plugins.go
    ├── your_mono_plugin
    │   ├── Dockerfile
    │   ├── README.md
    │   ├── cmd
    │   │   ├── root.go
    │   │   ├── server.go
    │   │   └── version.go
    │   ├── etc
    │   │   └── etc.yaml
    │   ├── internal
    │   │   ├── config
    │   │   │   └── config.go
    │   │   ├── custom
    │   │   │   └── custom.go
    │   │   ├── handler
    │   │   │   └── routes.go
    │   │   ├── middleware
    │   │   │   ├── middleware.go
    │   │   │   ├── response.go
    │   │   │   └── validator.go
    │   │   └── svc
    │   │       ├── config.go
    │   │       ├── middleware.go
    │   │       └── servicecontext.go
    │   ├── main.go
    │   └── serverless
    │       └── serverless.go
    └── your_plugin
        ├── Dockerfile
        ├── README.md
        ├── cmd
        │   ├── root.go
        │   ├── server.go
        │   └── version.go
        ├── etc
        │   └── etc.yaml
        ├── go.mod
        ├── internal
        │   ├── config
        │   │   └── config.go
        │   ├── custom
        │   │   └── custom.go
        │   ├── handler
        │   │   └── routes.go
        │   ├── middleware
        │   │   ├── middleware.go
        │   │   ├── response.go
        │   │   └── validator.go
        │   └── svc
        │       ├── config.go
        │       ├── middleware.go
        │       └── servicecontext.go
        ├── main.go
        └── serverless
            └── serverless.go
```
