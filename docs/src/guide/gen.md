---
title: 生成代码
icon: code
order: 4
---

jzero gen 根据 daemon/api 和 daemon/proto 文件生成代码. 生成代码的逻辑是调用 goctl 工具完成.

jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心.

```shell
jzero gen
```

执行命令后的代码结构为:

```shell
$ tree
.
├── cmd
│   ├── daemon.go
│   └── root.go
├── config.toml
├── daemon
│   ├── api
│   │   ├── app5.api
│   │   ├── file.api
│   │   └── hello.api
│   ├── daemon.go
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── handler
│   │   │   ├── file
│   │   │   │   ├── downloadhandler.go
│   │   │   │   └── uploadhandler.go
│   │   │   ├── hello
│   │   │   │   ├── helloparamhandler.go
│   │   │   │   ├── hellopathhandler.go
│   │   │   │   └── helloposthandler.go
│   │   │   ├── myhandler.go
│   │   │   ├── myroutes.go
│   │   │   └── routes.go
│   │   ├── logic
│   │   │   ├── credential
│   │   │   │   └── credentialversionlogic.go
│   │   │   ├── credentialv2
│   │   │   │   └── credentialversionlogic.go
│   │   │   ├── file
│   │   │   │   ├── downloadlogic.go
│   │   │   │   └── uploadlogic.go
│   │   │   ├── hello
│   │   │   │   ├── helloparamlogic.go
│   │   │   │   ├── hellopathlogic.go
│   │   │   │   └── hellopostlogic.go
│   │   │   ├── machine
│   │   │   │   └── machineversionlogic.go
│   │   │   └── machinev2
│   │   │       └── machineversionlogic.go
│   │   ├── server
│   │   │   ├── credential
│   │   │   │   └── credentialserver.go
│   │   │   ├── credentialv2
│   │   │   │   └── credentialv2server.go
│   │   │   ├── machine
│   │   │   │   └── machineserver.go
│   │   │   └── machinev2
│   │   │       └── machinev2server.go
│   │   ├── svc
│   │   │   └── servicecontext.go
│   │   └── types
│   │       └── types.go
│   ├── pb
│   │   ├── credentialpb
│   │   │   ├── credential.pb.go
│   │   │   └── credential_grpc.pb.go
│   │   └── machinepb
│   │       ├── machine.pb.go
│   │       └── machine_grpc.pb.go
│   └── proto
│       ├── credential.proto
│       └── machine.proto
├── go.mod
└── main.go

27 directories, 39 files
```