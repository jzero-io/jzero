# jzero

[![GitHub release](https://img.shields.io/github/release/jzero-io/jzero.svg?style=flat-square)](https://github.com/jzero-io/jzero/releases/latest)
[![Build Status](https://img.shields.io/github/actions/workflow/status/jzero-io/jzero/ci.yaml?branch=main&label=jzero-ci&logo=github&style=flat-square)](https://github.com/jzero-io/jzero/actions?query=workflow%3Ajzero-ci)
[![Go Report Card](https://goreportcard.com/badge/github.com/jzero-io/jzero?style=flat-square)](https://goreportcard.com/report/github.com/jzero-io/jzero)

<p align="center">
<img align="center" width="150px" src="https://oss.jaronnie.com/jzero.jpg">
</p>

基于 go-zero 框架定制的企业级后端代码框架 jzero

## Features

* 企业级代码规范
* grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要
* 集成命令行框架 cobra, 轻松编写具备生产可用的命令行工具
* 支持多 proto 多 service, 减少开发耦合性
* 不修改源码, 完全同步 go-zero 新特性
* 一键创建项目, 快速拓展新业务, 减少心理负担
* 一键生成服务端代码, 数据库代码, 客户端 sdk, 大大提高开发测试效率
* 支持自定义模板, 基于模板新建项目和生成代码
* 支持流量治理, 减少线上风险
* 支持链路追踪, 监控等, 快速定位问题
* 所有工具链跨平台支持

## [Quick Start](https://jzero.jaronnie.com/#快速开始)

![2024-04-30_10-10-52](https://oss.jaronnie.com/2024-04-30_10-10-52.gif)

### new project with grpc, gateway, api, cli

```shell
jzero new simple
```

### new project with only grpc

```shell
jzero new simplerpc --branch rpc
```

### new project with only api

```shell
jzero new simpleapi --branch api
```

### new project with only cli

```shell
jzero new simplecli --branch cli
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/jzero-io/jzero.svg)](https://starchart.cc/jzero-io/jzero)
