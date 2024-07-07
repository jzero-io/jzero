---
title: prometheus 配置
icon: vscode-icons:file-type-prometheus
star: true
order: 2
category: 配置
tag:
  - Guide
---

## go-zero grpc + grpc-gateway 项目

修改 etc/etc.yaml 添加一下内容

```shell
Zrpc:
  DevServer:
    Enabled: true
Gateway:
  DevServer:
    Enabled: true
```


## go-zero api 项目

修改 etc/etc.yaml 添加一下内容

```yaml
Rest:
  DevServer:
    Enabled: true
```
