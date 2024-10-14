# jzero

**Free your hands and have more time to play games**

[![Build Status](https://img.shields.io/github/actions/workflow/status/jzero-io/jzero/ci.yaml?branch=main&label=jzero-ci&logo=github&style=flat-square)](https://github.com/jzero-io/jzero/actions?query=workflow%3Ajzero-ci)
[![GitHub release](https://img.shields.io/github/release/jzero-io/jzero.svg?style=flat-square)](https://github.com/jzero-io/jzero/releases/latest)
[![GitHub package version](https://img.shields.io/github/v/release/jzero-io/jzero?include_prereleases&sort=semver&label=Docker%20Image%20version)](https://github.com/jzero-io/jzero/pkgs/container/jzero)
[![GitHub package version](https://img.shields.io/github/v/release/jzero-io/jzero-action?include_prereleases&sort=semver&label=Jzero%20Action%20Version)](https://github.com/marketplace/actions/jzero-action)

<p align="center">
<img align="center" width="150px" src="https://oss.jaronnie.com/jzero.jpg">
</p>

[中文](README.md) | ENGLISH

## Introduction

Create a project through a template and generate Server/Client/Model code based on proto/api/sql files.

It has the following features:

* Based on the [go-zero](https://go-zero.dev) framework but not limited to the go-zero framework. In theory, it can access any framework based on template features

* Optimize the existing pain points of the go-zero framework and expand new features, fully compatible with the go-zero framework

* Generate code through minimalist instructions based on configuration files

* Maintain commonly used development templates and generate projects that meet enterprise-level code specifications with one click

For more details please see: https://jzero.jaronnie.com/en. The English documentation is incomplete and is still being supplemented

## Install

```shell
go install github.com/jzero-io/jzero@latest
# check tools
jzero check
```

### docker

```shell
docker pull ghcr.io/jzero-io/jzero:latest
```

## Quick start

```shell
# new project
jzero new your_project
# generate server code
cd your_project
jzero gen
# download dependencies
go mod tidy
# generate swagger json
jzero gen swagger
# generate http sdk
jzero gen sdk
# generate zrpcclient
jzero gen zrpcclient
# run server
go run main.go server
```

### docker

```shell
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# download dependencies
go mod tidy
# generate swagger json
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# generate http sdk
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen sdk
# generate zrpcclient
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen zrpcclient
# 运行服务端
go run main.go server
```

For more examples code please see: https://github.com/jzero-io/examples

Project Practice please see: https://jzero.jaronnie.com/project

## Roadmap

please see: https://jzero.jaronnie.com/roadmap

## Contributors

[CONTRIBUTING](CONTRIBUTING.md)

<div>
  <a href="https://github.com/jzero-io/jzero/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=jzero-io/jzero" />
  </a>
</div>

## Acknowledgements

This project is supported by JetBrains Open Source development License.

[![Jetbrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://www.jetbrains.com/?from=jzero)

## Donate

[![Sponsor](https://img.shields.io/badge/Sponsor-%E2%9D%A4-red?label=Sponsor-WePay)](https://oss.jaronnie.com/2021723027876_.pic.jpg)
[![Sponsor](https://img.shields.io/badge/Sponsor-%E2%9D%A4-red?label=Sponsor-AliPay)](https://oss.jaronnie.com/2031723027877_.pic.jpg)

## Stargazers over time

[![Star History Chart](https://api.star-history.com/svg?repos=jzero-io/jzero&type=Date)](https://star-history.com/#jzero-io/jzero&Date)

## Contact me

<p align="center">
<img align="left" width="250px" height="250px" src="https://oss.jaronnie.com/weixin2.jpg">
</p>
