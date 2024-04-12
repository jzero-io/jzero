---
title: 指南
icon: lightbulb
---

## 安装

```shell
go install github.com/jaronnie/jzero@latest
# 初始化
jzero init
# 启动服务
jzero daemon
```

## 开发

```shell
# gencode
task gencode

# run
task run

# test
# unix
curl -s --unix-socket ./jzero.sock http://localhost:8001/api/v1.0/credential/version
# gateway
http://localhost:8001/api/v1.0/credential/version
# grpc
grpcurl -plaintext localhost:8000 credentialpb.credential/CredentialVersion
```