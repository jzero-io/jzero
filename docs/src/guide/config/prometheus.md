---
title: prometheus 配置
icon: gears
star: true
order: 2
category: 配置
tag:
  - Guide
---

修改 etc/etc.yaml 添加一下内容

```yaml
Zrpc:
  DevServer:
    Enabled: true
Gateway:
  DevServer:
    Enabled: true
```