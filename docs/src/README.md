---
home: false
icon: home
title: 首页
---

基于 go-zero 框架项目的代码设计

## 技术栈

* cobra 实现命令行管理
* go-zero 提供 grpc 和 http 请求等

## 特性

* 支持将 grpc 通过 gateway 转化为 http 请求, 并支持自定义 http 请求
* 同时支持在项目中使用 grpc, grpc-gateway, api
* 支持监听 unix socket
* 支持多 proto 多 service(多人开发友好)
* 一键创建项目(jzero new)
* 一键生成各种代码(jzero gen)

## 安装

```shell
go install github.com/jaronnie/jzero@latest
# 初始化
jzero init
# 启动样例服务
jzero daemon
```

## 一键创建项目

```shell
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest
# 一键安装相关工具
goctl env check --install --verbose --force
# 安装 jzero
go install github.com/jaronnie/jzero@latest
# 一键创建项目
jzero new --module=github.com/jaronnie/app1 --dir=./app1 --app=app1
cd app1
# 一键生成代码
jzero gen
# 下载依赖
go mod tidy
# 启动项目
go run main.go daemon --config config.toml
```

## 开发

```shell
# gencode
jzero gen

# run
go run main.go daemon --config config.toml

# test
# gateway
curl http://localhost:8001/api/v1.0/credential/version
# grpc
grpcurl -plaintext localhost:8000 credentialpb.credential/CredentialVersion
# api
curl http://localhost:8001/api/v1/hello/me
```

