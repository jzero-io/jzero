---
title: 链路追踪配置
icon: tabler:http-trace
star: true
order: 5
category: 配置
tag:
  - Guide
---

## rest

修改 etc/etc.yaml 添加一下内容

```yaml
rest:
  telemetry:
    name: "your_project-rpc"
    endpoint: "http://jaeger:14268/api/traces"
    sampler: 1.0
    batcher: "jaeger"
```

## zrpc

修改 etc/etc.yaml, 增加以下配置

```yaml
zrpc:
  telemetry:
    name: "your_project-rpc"
    endpoint: "http://jaeger:14268/api/traces"
    sampler: 1.0
    batcher: "jaeger"
```    

## gateway

修改 etc/etc.yaml 添加一下内容

```yaml
zrpc:
  telemetry:
    name: "your_project-rpc"
    endpoint: "http://jaeger:14268/api/traces"
    sampler: 1.0
    batcher: "jaeger"

gateway:
  telemetry:
    name: "your_project-gw"
    endpoint: "http://jaeger:14268/api/traces"
    sampler: 1.0
    batcher: "jaeger"
```