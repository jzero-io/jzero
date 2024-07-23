---
title: 概念介绍
icon: icon-park:concept-sharing
order: 1
---

基于 go-zero 开发的低代码脚手架 jzero, 旨在通过更少的命令完成更多的事情. 

该项目可一键创建项目, 并支持不同的使用场景, 如 grpc 项目, grpc + gateway 项目, go-zero api 项目以及命令行项目等. 通过项目的可描述文件(如 proto, api, sql 等)一键生成服务端代码, 客户端代码. 

并扩展了模板功能, 实现将任意代码转换为模板, 从而一键创建项目.

jzero 基于 go-zero 原生低代码脚手架 goctl 进行二次封装, 使用起来更加方便, 加强了项目规范, 并优化了以下特性:

* api 场景
  * 支持 types 文件分组(原生 goctl 将所有 api 文件生成的 types 放到单文件 types.go 中, 导致该文件爆炸)
* rpc 场景
  * 支持多个 proto(原生 goctl 仅支持单 proto, 多人开发下不够友好)
  * 默认支持 proto message 的字段校验, 且支持自定义错误信息
* gateway 场景
  * 默认可新增 rpc + gateway 的项目
  * 新增接口版本控制特性, 默认为 v1, 可一键初始化 v2, v3等版本的接口, 无需任何配置


