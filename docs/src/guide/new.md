---
title: 新建项目
icon: clone
order: 3
---

## 新建默认项目

其包含:

* grpc
* grpc-gateway
* api

::: code-tabs#shell

@tab jzero

```bash
jzero new app1
```

@tab Docker

```bash
docker run --rm -v ${PWD}/app1:/app/app1 jaronnie/jzero:latest new app1
```
:::

## 新建仅有 api 的项目

::: code-tabs#shell

@tab jzero

```bash
jzero new app1 --branch api
```

@tab Docker

```bash
docker run --rm -v ${PWD}/app1:/app/app1 jaronnie/jzero:latest new app1 --branch api
```
:::

## 新建仅有 rpc 的项目

::: code-tabs#shell

@tab jzero

```bash
jzero new app1 --branch rpc
```

@tab Docker

```bash
docker run --rm -v ${PWD}/app1:/app/app1 jaronnie/jzero:latest new app1 --branch rpc
```
:::

## 新建仅有 cli 的项目

::: code-tabs#shell

@tab jzero

```bash
jzero new app1 --branch cli
```

@tab Docker

```bash
docker run --rm -v ${PWD}/app1:/app/app1 jaronnie/jzero:latest new app1 --branch cli
```
:::

all flags:

* module 设置 go module
* dir 设置生成的项目路径
* home 设置本地 templates 路径
* remote 设置远程 templates repo
* branch 远程 templates repo branch


