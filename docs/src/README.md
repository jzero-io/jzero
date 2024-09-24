---
home: false
icon: fluent:home-heart-20-filled
title: 首页
---

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/jzero.jpg" style="width: 33%;" alt=""/>
</div>

## 特性

基于 [go-zero](https://go-zero.dev) 开发的低代码脚手架 jzero, 旨在通过更少的命令完成更多的事情, 但 jzero 不仅仅局限于 go-zero 框架, 理论上通过模板功能可以支持任意框架, 这取决于你如何使用它.

该项目可一键创建项目, 并支持不同的使用场景, 如 grpc 项目, gateway 项目, api 项目以及命令行项目等. 通过项目的可描述文件(如 proto, api, sql 等)一键生成服务端代码, 客户端代码和数据库代码.

jzero 具备以下特点:

* 极简命令, 通过配置文件 .jzero.yaml 控制不同命令的参数
* 具备不同场景下的开发模板, 具备快速复制项目的能力
* 优化 go-zero 已有的痛点并扩展新特性, 并完全兼容
* 模板特性支持新增任意文件, 基于模板特性理论上可以支持任意框架

在以下场景优化点:

* api 场景
    * 支持 types 文件分组(原生 goctl 将所有 api 文件生成的 types 放到单文件 types.go 中, 导致该文件爆炸)
    * 编写多个 api 文件, 无需显示的编写一个 main.api 文件, 框架自动处理好
    * 默认集成 `https://github.com/go-playground/validator` 校验框架
    * 支持重新生成 handler 文件, 并支持不同场景(有输入输出, 有输入没输出, 没输入有输出, 没输入没输出), 无需再手动维护 handler 代码, 提升开发过程中的体验
    * 支持自动修改 logic 文件函数的入参和出参, 当 api 文件修改后, 自动适配修改, 提升开发过程中的体验
* rpc 场景
    * 支持多个 proto, 自动注册, 无需手动编写
    * 默认支持 proto message 的字段校验, 且支持自定义错误信息
    * 默认支持通过 proto 新增拦截器, 可以设定某个 method, 也可以设定整个 service
* gateway 场景
    * 默认可新增 rpc + gateway 的项目
    * 新增接口版本控制特性, 默认为 v1, 可一键初始化 v2, v3等版本的接口, 无需任何配置
    * 默认支持通过 proto 新增拦截器和 http 中间件, 可以设定某个 method, 也可以设定整个 service
* 数据库场景
  * 将原生 sql 替换成 sqlbuilder, 从而可以更好的支持不同的数据库类型
  * 扩展新的抽象方法提升开发效率, 不再是简单的增删改查, 逐步扩展, 拥有类似 orm 的能力
* 客户端场景:
  * 通过 api/proto 文件自动生成 swagger json, 并内置 swagger ui 
  * 通过 api/proto 文件自动生成 http client
  * 通过 proto 文件自动生成 zrpc client

## 快速开始

:::tip Windows 用户请在 powershell 下执行所有指令
:::

::: code-tabs#shell

@tab Docker

```bash
# 一键创建项目
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project 
# 一键生成代码
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# 下载依赖
go mod tidy
# 启动项目
go run main.go server
```

@tab jzero

```bash
# 安装 jzero
go install github.com/jzero-io/jzero@latest
# 一键安装所需的工具
jzero check
# 一键创建项目
jzero new your_project
cd your_project
# 一键生成代码
jzero gen
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
```
:::

