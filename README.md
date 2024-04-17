# jzero

基于 go-zero 框架定制的企业级后端代码框架

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/jzero.jpg" style="width: 33%;" alt=""/>
</div>

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
* 一键编译各种平台的二进制并打包成镜像

## 快速开始

[详细文档点击这里](https://jzero.jaronnie.com)

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

## 验证

```shell
# test
# gateway
curl http://localhost:8001/api/v1.0/credential/version
# grpc
grpcurl -plaintext localhost:8000 credentialpb.credential/CredentialVersion
# api
curl http://localhost:8001/api/v1/hello/me
```

## RoadMap

- [x] Fix static embed files
- [x] Add go-zero api feature
- [x] Support multi proto, multi service
- [x] Support api ~~multi api,~~ multi service
- [x] Warp rpc and api Response
- [x] Support jzero gen
- [x] Support jzero new
- [x] jzero gen support auto register service to zrpc server and update gateway.UpStream.0.ProtoSets

