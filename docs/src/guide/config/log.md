---
title: 日志配置
icon: gears
star: true
order: 1
category: 配置
tag:
  - Guide
---

::: tip jzero version > v0.8.0 支持在日志模式为 file 或者 volume 的情况下仍然输出到控制台
默认 LogToConsole 值为 true. 如果需要关闭, 可以设置为 false.
:::

```toml
[Log]
ServiceName = "app1"
Level = "info"
Mode = "file"
encoding = "plain"
KeepDays = 30
MaxBackups = 7
MaxSize = 50
Rotation = "size"

[App1]
LogToConsole = true
```

默认配置下日志最大占用空间: 2G

计算规则如下: 

logs 文件夹一共 5 个文件. 每个文件最大占用 50MB, 最多备份 7 个. 即 50MB * 8 * 5


