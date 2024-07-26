---
title: prometheus 配置
icon: vscode-icons:file-type-prometheus
star: true
order: 2
category: 配置
tag:
  - Guide
---

## rest

修改 etc/etc.yaml, 增加以下配置

```yaml
rest:
  devServer:
    enabled: true
```

## zrpc

修改 etc/etc.yaml, 增加以下配置

```shell
zrpc:
  devServer:
    enabled: true
```

## gateway

修改 etc/etc.yaml, 增加以下配置

```shell
zrpc:
  devServer:
    enabled: true
gateway:
  devServer:
    enabled: true
```