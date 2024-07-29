---
title: 日志配置
icon: mdi:math-log
star: true
order: 1
category: 配置
tag:
  - Guide
---

修改 etc/etc.yaml, 增加以下配置

```yaml
log:
  keepDays: 30
  level: info
  maxBackups: 7
  maxSize: 50
  mode: file
  rotation: size
  serviceName: your_project
  encoding: plain
```

默认配置下日志最大占用空间: 2G

计算规则如下: 

logs 文件夹一共 5 个文件. 每个文件最大占用 50MB, 最多备份 7 个. 即 50MB * 8 * 5