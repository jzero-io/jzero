---
title: 链路追踪配置
icon: tabler:http-trace
star: true
order: 3
category: 配置
tag:
  - Guide
---

## Rest

修改 etc/etc.yaml 添加一下内容

```yaml
Rest:
  Telemetry:
    Name: "your_project-rpc"
    Endpoint: "http://jaeger:14268/api/traces"
    Sampler: 1.0
    Batcher: "jaeger"
```

## Gateway

修改 etc/etc.yaml 添加一下内容

```yaml
Zrpc:
  Telemetry:
    Name: "your_project-rpc"
    Endpoint: "http://jaeger:14268/api/traces"
    Sampler: 1.0
    Batcher: "jaeger"

Gateway:
  Telemetry:
    Name: "your_project-gw"
    Endpoint: "http://jaeger:14268/api/traces"
    Sampler: 1.0
    Batcher: "jaeger"
```