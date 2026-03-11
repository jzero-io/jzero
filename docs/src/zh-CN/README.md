---
home: false
icon: fluent:home-heart-20-filled
title: 首页
---

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/jzero.svg" style="width: 33%;" alt=""/>
</div>

## 简介

基于 [go-zero框架](https://github.com/zeromicro/go-zero) 以及 [go-zero/goctl工具](https://github.com/zeromicro/go-zero/tree/master/tools/goctl) 开发的 [jzero](https://github.com/jzero-io/jzero) 框架, 可一键初始化 api/gateway/rpc 项目。

基于可描述文件(**api/proto/sql**)自动生成**服务端和客户端**框架代码, 基于内置的 jzero-skills 让 AI 生成符合最佳实践的业务逻辑代码，降低开发心智, 解放双手!

具备以下特点:

* 支持通过**配置文件/命令行参数/环境变量**组合的方式灵活控制 jzero 的各项配置, 极简指令生成代码, ai 友好
* 支持基于 **git 对改动文件**生成代码, 支持对**指定描述文件**生成代码或**忽略指定描述文件**生成代码, 提升大型项目代码生成效率
* 内置常用开发模板并增强模板特性, 支持**自定义模板**, 构建专属企业内部代码模板, 极大降低开发成本
* 支持**插件化架构**, 功能模块可作为独立插件动态加载, 支持插件创建、编译和卸载, 完美适配团队协作和模块解耦

更多详情请参阅：[https://docs.jzero.io](https://docs.jzero.io)

## 设计理念

* **开发体验**: 提供简单好用的一站式生产可用的解决方案, 提升开发体验感
* **模板驱动**: 所有代码生成均基于模板渲染, 默认生成即最佳实践, 且支持自定义模板内容
* **生态兼容**: 不修改 go-zero 和 go-zero/goctl, 保持生态兼容, 同时解决已有的痛点问题并扩展新的功能
* **团队开发**: 通过模块**分层**, **插件**设计, 团队开发友好
* **接口设计**: 不依赖特定数据库/缓存/配置中心等基础设施, 根据实际需求自由选择

## 快速开始

::: code-tabs#shell
@tab jzero cli

```bash
# 安装 jzero
go install github.com/jzero-io/jzero/cmd/jzero@latest
# 一键安装所需的工具
jzero check
# 一键创建项目
jzero new your_project
cd your_project
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```

@tab jzero Docker

```bash
# 一键创建项目
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```
:::



