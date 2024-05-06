---
home: false
icon: home
title: 首页
---

基于 [go-zero](https://go-zero.dev) 框架定制的企业级后端代码框架

::: tip 目前还在定制规范中, 不能保证兼容性. 待 v1.0 后保证兼容性
:::

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/jzero.jpg" style="width: 33%;" alt=""/>
</div>

## 特性

* 企业级代码规范
* grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要
* 支持监听 unix socket 本地通信
* 支持多 proto 多 service, 减少开发耦合性
* 不修改源码, 完全同步 go-zero 新特性
* 一键创建项目, 快速拓展新业务, 减少心理负担
* 一键生成服务端代码, 数据库代码, 客户端 sdk, 大大提高开发测试效率
* 支持自定义模板, 基于模板新建项目和生成代码
* 支持流量治理, 减少线上风险
* 支持链路追踪, 监控等, 快速定位问题
* 所有工具链跨平台支持

## 快速开始

![2024-04-30_10-10-52](https://oss.jaronnie.com/2024-04-30_10-10-52.gif)

:::tip Windows 用户请在 powershell 下执行所有指令
:::

::: code-tabs#shell

@tab Docker(amd64)

```bash
# 一键创建项目
docker run --rm -v ${PWD}/app1:/app/app1 jaronnie/jzero:latest new --module=github.com/jaronnie/app1 --dir=./app1 --app=app1
  
cd app1 

# 一键生成代码
docker run --rm -v ${PWD}:/app/app1 jaronnie/jzero:latest gen -w app1

# 下载依赖
go mod tidy
# 启动项目
go run main.go daemon
```

@tab Docker(arm64)

```bash
# 一键创建项目
docker run --rm -v ${PWD}/app1:/app/app1 jaronnie/jzero:latest-arm64 new --module=github.com/jaronnie/app1 --dir=./app1 --app=app1
  
cd app1  

# 一键生成代码
docker run --rm -v ${PWD}:/app/app1 jaronnie/jzero:latest-arm64 gen -w app1

# 下载依赖
go mod tidy
# 启动项目
go run main.go daemon
```

@tab jzero

```bash
# 安装 jzero
go install github.com/jzero-io/jzero@latest
# 一键安装所需的工具
jzero check
# 一键创建项目
jzero new --module=github.com/jaronnie/app1 --dir=./app1 --app=app1
cd app1
# 一键生成代码
jzero gen
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go daemon
```
:::

## 验证

```shell
# test
# gateway
curl http://localhost:8001/api/v1/credential/version
# grpc
grpcurl -plaintext localhost:8000 credentialpb.credential/CredentialVersion
# api
curl http://localhost:8001/api/v1/hello/me
```

