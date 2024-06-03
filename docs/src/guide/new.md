---
title: 新建项目
icon: clone
order: 3
---

## go-zero with grpc + gateway + api

::: code-tabs#shell

@tab jzero

```bash
jzero new quickstart
```

@tab Docker

```bash
docker run --rm -v ${PWD}/quickstart:/app/quickstart jaronnie/jzero:latest new quickstart
```
:::

## go-zero with grpc + gateway

::: code-tabs#shell

@tab jzero

```bash
jzero new simplegateway --branch gateway
```

@tab Docker

```bash
docker run --rm -v ${PWD}/simplegateway:/app/simplegateway jaronnie/jzero:latest new simplegateway --branch gateway
```
:::

## go-zero with api

::: code-tabs#shell

@tab jzero

```bash
jzero new simpleapi --branch api
```

@tab Docker

```bash
docker run --rm -v ${PWD}/simpleapi:/app/simpleapi jaronnie/jzero:latest new simpleapi --branch api
```
:::

## go-zero with zrpc

::: code-tabs#shell

@tab jzero

```bash
jzero new simplerpc --branch rpc
```

@tab Docker

```bash
docker run --rm -v ${PWD}/simplerpc:/app/simplerpc jaronnie/jzero:latest new simplerpc --branch rpc
```
:::

## with cobra cli project

::: code-tabs#shell

@tab jzero

```bash
jzero new simplecli --branch cli
```

@tab Docker

```bash
docker run --rm -v ${PWD}/simplecli:/app/simplecli jaronnie/jzero:latest new simplecli --branch cli
```
:::


