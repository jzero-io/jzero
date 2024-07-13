---
title: prometheus 配置
icon: vscode-icons:file-type-prometheus
star: true
order: 2
category: 配置
tag:
  - Guide
---

## Rest

修改 etc/etc.yaml 添加一下内容

```yaml
rest:
  devServer:
    enabled: true
```


## Gateway

修改 etc/etc.yaml 添加一下内容

```shell
zrpc:
  devServer:
    enabled: true
gateway:
  devServer:
    enabled: true
```