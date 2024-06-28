# jzero

[![GitHub release](https://img.shields.io/github/release/jzero-io/jzero.svg?style=flat-square)](https://github.com/jzero-io/jzero/releases/latest)
[![Build Status](https://img.shields.io/github/actions/workflow/status/jzero-io/jzero/ci.yaml?branch=main&label=jzero-ci&logo=github&style=flat-square)](https://github.com/jzero-io/jzero/actions?query=workflow%3Ajzero-ci)
[![Go Report Card](https://goreportcard.com/badge/github.com/jzero-io/jzero?style=flat-square)](https://goreportcard.com/report/github.com/jzero-io/jzero)

<p align="center">
<img align="center" width="150px" src="https://oss.jaronnie.com/jzero.jpg">
</p>

可支持任意框架的脚手架 jzero, 默认支持 go-zero

## Features

* 企业级代码规范, 多人开发友好
* 支持自定义模板, 基于模板新建项目和生成代码, 默认支持多场景开发模板
  * grpc + grpc gateway + api
  * api
  * grpc
  * grpc + grpc gateway
  * cli
* 扩展 go-zero 功能，完全同步 go-zero 新特性
  * 保持与原生 goctl 生成的目录, 易迁移过来
  * 支持 types 分组
  * 支持删除冗余的 Logic, Handler, Server 等字样
* 一键创建项目
* 一键生成服务端代码, 数据库代码, 客户端 sdk

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

## Contact me

* Weixin: jaronnie
* Email: jaron@jaronnie.com
