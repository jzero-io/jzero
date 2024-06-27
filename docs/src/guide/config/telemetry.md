---
title: 链路追踪配置
icon: gears
star: true
order: 3
category: 配置
tag:
  - Guide
---

修改 etc/etc.yaml 添加一下内容

```yaml
Zrpc:
  Telemetry:
    Name: "app1-rpc"
    Endpoint: "http://jaeger:14268/api/traces"
    Sampler: 1.0
    Batcher: "jaeger"

Gateway:
  Telemetry:
    Name: "app1-gw"
    Endpoint: "http://jaeger:14268/api/traces"
    Sampler: 1.0
    Batcher: "jaeger"
```