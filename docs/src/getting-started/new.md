---
title: Project Initialization
icon: mdi:new-box
order: 3
---

## Template Introduction

A template is a predefined set of code structures that provides the basic architecture and engineering standards for a project.

Templates help you quickly start initializing a project without writing code from scratch.

## Template Types

jzero provides the following types of templates to meet various scenarios:

* Built-in template(frame): Built-in template providing core framework capabilities, supports optional features (database/cache)
* Path template(home): Specify a path as a template, usually placed inside a specific project to meet specific project needs
* Local template(local): Local global template located in ~/.jzero/templates/local folder
* Remote repository template(remote+branch): Can be used to build enterprise-specific remote template repositories

For detailed usage, see: [Template Guide](../guide/template.md)

## Initialize api project

::: code-tabs#shell

@tab jzero cli

```bash
jzero new your_project --frame api
cd your_project
# download dependencies
go mod tidy
# start server
go run main.go server
# visit swagger ui
http://localhost:8001/swagger
```

@tab jzero Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame api
cd your_project
# download dependencies
go mod tidy
# start server
go run main.go server
# visit swagger ui
http://localhost:8001/swagger
```
:::

## Initialize rpc project

::: code-tabs#shell

@tab jzero cli

```bash
jzero new your_project --frame rpc
cd your_project
# download dependencies
go mod tidy
# start server
go run main.go server
```

@tab jzero Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame rpc
cd your_project
# download dependencies
go mod tidy
# start server
go run main.go server
```
:::

## Initialize gateway project

:::important Supports both grpc/http interfaces
:::

::: code-tabs#shell

@tab jzero cli

```bash
jzero new your_project --frame gateway
cd your_project
# download dependencies
go mod tidy
# start server
go run main.go server
# visit swagger ui
http://localhost:8001/swagger
```

@tab jzero Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame gateway
cd your_project
# download dependencies
go mod tidy
# start server
go run main.go server
# visit swagger ui
http://localhost:8001/swagger
```
:::

## Optional features model/redis/model+redis

Based on optional features, provides a complete solution for using model/redis/model

```shell
# Use case: need to connect to relational database(model) with database cache(cache), redis
jzero new your_project --features model,cache,redis

# Use case: need to connect to relational database(model), redis
jzero new your_project --features model,redis

# Use case: need to connect to relational database(model) with database cache(cache)
jzero new your_project --features model,cache

# Use case: need to connect to relational database(model)
jzero new your_project --features model
```
