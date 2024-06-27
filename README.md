# jzero

[![GitHub release](https://img.shields.io/github/release/jzero-io/jzero.svg?style=flat-square)](https://github.com/jzero-io/jzero/releases/latest)
[![Build Status](https://img.shields.io/github/actions/workflow/status/jzero-io/jzero/ci.yaml?branch=main&label=jzero-ci&logo=github&style=flat-square)](https://github.com/jzero-io/jzero/actions?query=workflow%3Ajzero-ci)
[![Go Report Card](https://goreportcard.com/badge/github.com/jzero-io/jzero?style=flat-square)](https://goreportcard.com/report/github.com/jzero-io/jzero)

<p align="center">
<img align="center" width="150px" src="https://oss.jaronnie.com/jzero.jpg">
</p>

可支持任意框架的脚手架 jzero, 默认支持 go-zero

## Features

* 企业级代码规范
* grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要
* 集成命令行框架 cobra, 轻松编写具备生产可用的命令行工具
* 支持多 proto 多 service, 减少开发耦合性
* 不修改源码, 完全同步 go-zero 新特性
* 一键创建项目, 快速拓展新业务, 减少心理负担
* 一键生成服务端代码, 数据库代码, 客户端 sdk, 大大提高开发测试效率
* 支持自定义模板, 基于模板新建项目和生成代码

## [Quick Start](https://jzero.jaronnie.com/#快速开始)

```shell
go install github.com/jzero-io/jzero@latest
jzero check
jzero new your_project
cd your_project && jzero gen && go mod tidy
go run main.go server
```


## Stargazers over time

[![Stargazers over time](https://starchart.cc/jzero-io/jzero.svg)](https://starchart.cc/jzero-io/jzero)
