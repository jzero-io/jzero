---
title: 监听 unix sock 配置
icon: gears
star: true
order: 4
category: 配置
tag:
  - Guide
---

修改 config.yaml. 添加一下内容

::: tip 以创建项目时填写的 app 名称的首字母大写作为配置项, 可查看 config.yaml APP 查看值
:::

```yaml
App1:
  ListenOnUnixSocket: "./app1.sock"
```

```shell
Using config file: config.toml
2024-04-19T13:39:17.708+08:00    info   Starting dev http server at :6060       caller=devserver/server.go:71
Starting rpc server at 0.0.0.0:8000...
Starting gateway server at 0.0.0.0:8001...
Starting unix server at ./app1.sock...

```

## 验证

```shell
curl -s --unix-socket ./app1.sock http://localhost/api/v1/credential/version  
```