---
title: 限流配置
icon: mdi:car-speed-limiter
star: true
order: 6
category: 配置
tag:
  - Guide
---

## rest

修改 etc/etc.yaml, 增加以下配置, 设置最大 qps 100

```yaml
rest:
  maxConns: 100
```

## gateway

修改 etc/etc.yaml, 增加以下配置, 设置最大 qps 100

```yaml
gateway:
  maxConns: 100
```