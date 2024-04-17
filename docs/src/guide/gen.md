---
title: 生成代码
icon: code
order: 4
---

jzero gen 根据 daemon/api 和 daemon/proto 文件生成代码. 生成代码的逻辑是调用 goctl 工具完成.

jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心.

## 生成代码

```shell
cd app1
jzero gen
```

## 下载依赖

```shell
go mod tidy
```

## 运行项目

```shell
go run main.go daemon --config config.toml
```

## 测试接口

```shell
# gateway
curl http://localhost:8001/api/v1.0/credential/version
# grpc
grpcurl -plaintext localhost:8000 credentialpb.credential/CredentialVersion
# api
curl http://localhost:8001/api/v1/hello/me
```