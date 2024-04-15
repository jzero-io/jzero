---
title: 新建项目
icon: clone
order: 3
---

```shell
jzero new --module=github.com/jaronnie/app1 --dir=./app1 --app=app1
```

flag 解释:

* module 表示新建项目的 go module
* dir 表示创建的项目目录路径
* app 表示项目名

生成的代码结构:

```shell
$ tree                           
.
├── cmd
│   ├── daemon.go
│   └── root.go
├── config.toml
├── daemon
│   ├── api
│   │   ├── app1.api
│   │   ├── file.api
│   │   └── hello.api
│   ├── daemon.go
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   └── handler
│   │       ├── myhandler.go
│   │       └── myroutes.go
│   └── proto
│       ├── credential.proto
│       └── machine.proto
├── go.mod
└── main.go

8 directories, 14 files
```