---
title: 超时配置
icon: eos-icons:timeout
order: 0
star: true
category: 配置
tag:
  - Guide
---

## rest

修改 etc/etc.yaml, 增加以下配置

```yaml
rest:
  timeout: 10000 # 10s
```

## zrpc

修改 etc/etc.yaml, 增加以下配置

```yaml
zrpc:
  timeout: 10000 # 10s
```

## gateway

修改 etc/etc.yaml, 增加以下配置

```yaml
gateway:
  timeout: 10000 # 10s
  upstreams:
    - grpc:
        timeout: 10000 # 10s
zrpc:
  timeout: 10000 # 10s
```