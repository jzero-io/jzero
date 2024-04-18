---
title: 链路追踪配置
icon: gears
star: true
order: 3
category: 配置
tag:
  - Guide
---

修改 config.toml. 添加一下内容

```toml
[Telemetry]
Name = "your_app-gw"
Endpoint = "http://jaeger:14268/api/traces"
Sampler = 1.0
Batcher = "jaeger"

[Gateway.Telemetry]
Name = "your_app-gw"
Endpoint = "http://jaeger:14268/api/traces"
Sampler = 1.0
Batcher = "jaeger"
```