---
title: 限流配置
icon: mdi:car-speed-limiter
star: true
order: 5
category: 配置
tag:
  - Guide
---

## Rest

修改 etc/etc.yaml, 增加以下配置, 设置最大 qps 100

```yaml
Rest:
  MaxConns: 100
```

## Gateway

修改 etc/etc.yaml, 增加以下配置, 设置最大 qps 100

```yaml
Gateway:
  MaxConns: 100
```