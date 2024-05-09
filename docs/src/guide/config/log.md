---
title: 日志配置
icon: gears
star: true
order: 1
category: 配置
tag:
  - Guide
---

修改 config.yaml

```yaml
Log:
  KeepDays: 30
  Level: info
  MaxBackups: 7
  MaxSize: 50
  Mode: file
  LogToConsole: true # 将日志重定向到 console.
  Rotation: size
  ServiceName: app1
  encoding: plain
```

默认配置下日志最大占用空间: 2G

计算规则如下: 

logs 文件夹一共 5 个文件. 每个文件最大占用 50MB, 最多备份 7 个. 即 50MB * 8 * 5


