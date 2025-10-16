# jzero

**Free your hands and have more time to play games**

[![Build Status](https://img.shields.io/github/actions/workflow/status/jzero-io/jzero/ci.yaml?branch=main&label=jzero-ci&logo=github&style=flat-square)](https://github.com/jzero-io/jzero/actions?query=workflow%3Ajzero-ci)
[![GitHub release](https://img.shields.io/github/release/jzero-io/jzero.svg?style=flat-square)](https://github.com/jzero-io/jzero/releases/latest)
[![GitHub package version](https://img.shields.io/github/v/release/jzero-io/jzero?include_prereleases&sort=semver&label=Docker%20Image%20version)](https://github.com/jzero-io/jzero/pkgs/container/jzero)
[![GitHub package version](https://img.shields.io/github/v/release/jzero-io/jzero-action?include_prereleases&sort=semver&label=Jzero%20Action%20Version)](https://github.com/marketplace/actions/jzero-action)
[![Endpoint Badge](https://img.shields.io/endpoint?url=https%3A%2F%2Fvercel.jaronnie.com%2Fapi%2Fv1%2Fshields%2Fgithub%2Fjzero-io%2Fjzero%2Fpkgs%2Fcontainer%2Fjzero%2Fdownloads&label=image%20pulls)](https://vercel.jaronnie.com/api/v1/shields/github/jzero-io/jzero/pkgs/container/jzero/downloads)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/jzero-io/jzero)

<p align="center">
<img align="center" width="150px" src="https://oss.jaronnie.com/jzero.jpg">
</p>

![Static Badge](https://img.shields.io/badge/Latest_New_Feature-blue?style=for-the-badge)

* Deploy jzero applications on [Vercel](https://vercel.com), [Template repository branch](https://github.com/jzero-io/templates/tree/api-vercel), [Code example](https://github.com/jaronnie/go-serverless-vercel)
* Generate code based on git changes, greatly improving development experience for large projects
* [Admin management system](https://github.com/jzero-io/jzero-admin) based on jzero, [Demo 1 - deployed on Vercel](https://admin.jzero.io), [Demo 2 - deployed on Alibaba Cloud Function Compute](http://jzero-admin.jaronnie.com)
* [Serverless](https://docs.jzero.io/guide/serverless.html) plugin feature based on jzero, achieving multi-module decoupling and automatic dependency injection of third-party modules at compile time
* [Dynamic configuration feature](https://docs.jzero.io/guide/config/dynamic_conf.html), allowing dynamic modification of server configuration without restarting the server
* Implemented jzero mcp server, allowing jzero tools to be called in mcp client dialogs to generate code
* The same codebase can dynamically adapt to multiple database types

[中文](README.md) | ENGLISH

## Introduction

Create projects through templates and generate Server/Client/Model code based on proto/api/sql files.

It has the following features:
* Based on the [go-zero](https://go-zero.dev) framework but not limited to the go-zero framework. In theory, it can integrate with any framework based on template features
* Optimize existing pain points of the go-zero framework and extend new features, fully compatible with the go-zero framework
* Generate code through minimalist commands based on configuration files
* Generate code based on git changes only for modified files, greatly improving code generation efficiency for large projects
* Maintain commonly used development templates and generate projects that meet enterprise-level code specifications with one click
* All supporting toolchains are cross-platform, supporting Windows/Mac/Linux

For more details please see: https://docs.jzero.io

## Download

```shell
go install github.com/jzero-io/jzero/cmd/jzero@latest
# Check tools and download
jzero check
```

### docker

```shell
docker pull ghcr.io/jzero-io/jzero:latest
```

## Quick Start

```shell
# Create new project
jzero new your_project
# Generate server code
cd your_project
jzero gen
# Download dependencies
go mod tidy
# Generate swagger json
jzero gen swagger
# Generate http client sdk
jzero gen sdk
# Run server
go run main.go server
```

### docker

```shell
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# Download dependencies
go mod tidy
# Generate swagger json
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# Generate http client sdk
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen sdk
# Run server
go run main.go server
```

For more example code please see: https://github.com/jzero-io/examples

Project practice please see:
* api framework: https://docs.jzero.io/project/api.html
* gateway framework: https://docs.jzero.io/project/gateway.html

## Related Projects
* jzero-admin (Admin management system based on jzero): https://github.com/jzero-io/jzero-admin

## Roadmap

Please see: https://docs.jzero.io/roadmap/base.html

## Contributors

[Contributing](https://docs.jzero.io/guide/contribute.html)

<a href="https://openomy.app/github/jzero-io/jzero" target="_blank" style="display: block; width: 100%;" align="center">
  <img src="https://openomy.app/svg?repo=jzero-io/jzero&chart=bubble&latestMonth=3" target="_blank" alt="Contribution Leaderboard" style="display: block; width: 100%;" />
</a>

## Acknowledgements

This project is supported by JetBrains Open Source development License.

[![Jetbrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://www.jetbrains.com/?from=jzero)

## Donate

[![Sponsor](https://img.shields.io/badge/Sponsor-%E2%9D%A4-red?label=Sponsor-WePay)](https://oss.jaronnie.com/2021723027876_.pic.jpg)
[![Sponsor](https://img.shields.io/badge/Sponsor-%E2%9D%A4-red?label=Sponsor-AliPay)](https://oss.jaronnie.com/2031723027877_.pic.jpg)

## Stargazers over time

[![Star History Chart](https://api.star-history.com/svg?repos=jzero-io/jzero&type=Date)](https://star-history.com/#jzero-io/jzero&Date)

## Disclaimer

This project is for learning and communication purposes only. Please do not use it for illegal purposes. The author is not responsible for any consequences arising from the use of this project.

## Contact

<p align="center">
<img align="left" width="250px" height="250px" src="https://oss.jaronnie.com/weixin2.jpg">
</p>
