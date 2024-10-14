# jzero

**解放你的双手有更多的时间去玩游戏**

[![Build Status](https://img.shields.io/github/actions/workflow/status/jzero-io/jzero/ci.yaml?branch=main&label=jzero-ci&logo=github&style=flat-square)](https://github.com/jzero-io/jzero/actions?query=workflow%3Ajzero-ci)
[![GitHub release](https://img.shields.io/github/release/jzero-io/jzero.svg?style=flat-square)](https://github.com/jzero-io/jzero/releases/latest)
[![GitHub package version](https://img.shields.io/github/v/release/jzero-io/jzero?include_prereleases&sort=semver&label=Docker%20Image%20version)](https://github.com/jzero-io/jzero/pkgs/container/jzero)
[![GitHub package version](https://img.shields.io/github/v/release/jzero-io/jzero-action?include_prereleases&sort=semver&label=Jzero%20Action%20Version)](https://github.com/marketplace/actions/jzero-action)
[![Endpoint Badge](https://img.shields.io/endpoint?url=https%3A%2F%2Fvercel.jaronnie.com%2Fapi%2Fv1%2Fshields%2Fgithub%2Fjzero-io%2Fjzero%2Fpkgs%2Fcontainer%2Fjzero%2Fdownloads&label=image%20pulls)](https://vercel.jaronnie.com/api/v1/shields/github/jzero-io/jzero/pkgs/container/jzero/downloads)

<p align="center">
<img align="center" width="150px" src="https://oss.jaronnie.com/jzero.jpg">
</p>

![Static Badge](https://img.shields.io/badge/Latest_New_Feature-blue?style=for-the-badge)

* 将 jzero 应用部署在 [Vercel](https://vercel.com), [模板仓库分支](https://github.com/jzero-io/templates/tree/api-vercel), [代码示例](https://github.com/jaronnie/go-serverless-vercel)
* 基于 git change 生成代码, 极大提升大型项目开发体验

中文 | [ENGLISH](README-EN.md)

## 介绍

通过模板创建项目，并基于 proto/api/sql 文件生成 Server/Client/Model 代码。

具备以下特点:
* 基于 [go-zero](https://go-zero.dev) 框架但不局限于 go-zero 框架, 理论上可以基于模板特性接入任意框架
* 优化 go-zero 框架已有痛点, 并扩展新的特性, 完全兼容 go-zero 框架
* 基于配置文件, 通过极简指令生成代码
* 基于 git 仅对改动文件部分生成代码, 极大提升大型项目代码生成效率
* 维护常用开发模板, 一键生成符合企业级代码规范的项目
* 所有配套工具链跨平台使用, 支持 windows/mac/linux

更多详情请参阅：https://jzero.jaronnie.com

## 下载

```shell
go install github.com/jzero-io/jzero@latest
# 检查工具并下载
jzero check
```

### docker

```shell
docker pull ghcr.io/jzero-io/jzero:latest
```

## 快速开始

```shell
# 新建项目
jzero new your_project
# 生成服务端代码
cd your_project
jzero gen
# 下载依赖
go mod tidy
# 生成 swagger json
jzero gen swagger
# 生成 http 客户端 sdk
jzero gen sdk
# 生成 zrpc 客户端 sdk
jzero gen zrpcclient
# 运行服务端
go run main.go server
```

### docker

```shell
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# 下载依赖
go mod tidy
# 生成 swagger json
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# 生成 http 客户端 sdk
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen sdk
# 生成 zrpc 客户端 sdk
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen zrpcclient
# 运行服务端
go run main.go server
```

更多示例代码请参阅: https://github.com/jzero-io/examples

项目实战请参阅: https://jzero.jaronnie.com/project

## 路线图

请参阅: https://jzero.jaronnie.com/roadmap

## 贡献者

[贡献](https://jzero.jaronnie.com/guide/contribute)

<div>
  <a href="https://github.com/jzero-io/jzero/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=jzero-io/jzero" />
  </a>
</div>

## 致谢

该项目由 JetBrains 开源开发许可证支持。

[![Jetbrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://www.jetbrains.com/?from=jzero)

## 捐赠

[![Sponsor](https://img.shields.io/badge/Sponsor-%E2%9D%A4-red?label=Sponsor-WePay)](https://oss.jaronnie.com/2021723027876_.pic.jpg)
[![Sponsor](https://img.shields.io/badge/Sponsor-%E2%9D%A4-red?label=Sponsor-AliPay)](https://oss.jaronnie.com/2031723027877_.pic.jpg)

## Stargazers over time

[![Star History Chart](https://api.star-history.com/svg?repos=jzero-io/jzero&type=Date)](https://star-history.com/#jzero-io/jzero&Date)

## 联系我

<p align="center">
<img align="left" width="250px" height="250px" src="https://oss.jaronnie.com/weixin2.jpg">
</p>
