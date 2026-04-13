# jzero

[![Build Status](https://img.shields.io/github/actions/workflow/status/jzero-io/jzero/ci.yaml?branch=main&label=jzero-ci&logo=github&style=flat-square)](https://github.com/jzero-io/jzero/actions?query=workflow%3Ajzero-ci)
[![GitHub release](https://img.shields.io/github/release/jzero-io/jzero.svg?style=flat-square)](https://github.com/jzero-io/jzero/releases/latest)
[![GitHub package version](https://img.shields.io/github/v/release/jzero-io/jzero-action?include_prereleases&sort=semver&label=Jzero%20Action%20Version)](https://github.com/marketplace/actions/jzero-action)
[![Endpoint Badge](https://img.shields.io/endpoint?url=https%3A%2F%2Fvercel.jaronnie.com%2Fapi%2Fv1%2Fshields%2Fgithub%2Fjzero-io%2Fjzero%2Fpkgs%2Fcontainer%2Fjzero%2Fdownloads&label=image%20pulls)](https://vercel.jaronnie.com/api/v1/shields/github/jzero-io/jzero/pkgs/container/jzero/downloads)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/jzero-io/jzero)

<p align="center">
<img align="center" width="150px" src="https://oss.jaronnie.com/jzero.svg">
<img align="center" width="300px" src="https://github.com/user-attachments/assets/44184df0-20ce-403d-ab38-74088915bc33">

</p>

**English** | **[简体中文](README.zh-CN.md)**

## Introduction

[jzero](https://github.com/jzero-io/jzero) is a framework developed based on the [go-zero framework](https://github.com/zeromicro/go-zero) and [go-zero/goctl tool](https://github.com/zeromicro/go-zero/tree/master/tools/goctl). It can initialize api/gateway/rpc projects with a single command.

Automatically generate **server and client** framework code based on descriptive files (**api/proto/sql**). With built-in jzero-skills, AI can generate business logic code that follows best practices, reducing development cognitive load and freeing your hands!

Key features:

* Flexible control of jzero configurations through **configuration files/command-line arguments/environment variables**, simple commands to generate code, AI-friendly
* Support generating code based on **git changed files**, support generating code for **specified descriptive files** or **ignoring specified descriptive files**, improving code generation efficiency for large projects
* Built-in common development templates with enhanced template features, support for **custom templates**, building exclusive enterprise internal code templates, greatly reducing development costs

For more details, please visit: [https://docs.jzero.io](https://docs.jzero.io)

## Design Philosophy

* **Developer Experience**: Provide a simple, easy-to-use, production-ready solution to enhance developer experience
* **Template Driven**: All code generation is based on template rendering, default generation follows best practices, and supports custom template content
* **Ecosystem Compatibility**: Does not modify go-zero and go-zero/goctl, maintains ecosystem compatibility, while solving existing pain points and extending new features
* **Team Development**: Team development friendly through module **layering** and **plugin** design
* **Interface Design**: No dependency on specific databases/caches/configuration centers and other infrastructure, choose freely according to actual needs

For more details, please visit: https://docs.jzero.io

## Quick Start

```shell
# Install jzero
go install github.com/jzero-io/jzero/cmd/jzero@latest
# One-click install required tools
jzero check
# One-click create project
jzero new your_project
cd your_project
# Download dependencies
go mod tidy
# Start server
go run main.go server
# Access swagger ui
http://localhost:8001/swagger
```

### docker

```shell
# One-click create project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project
# Download dependencies
go mod tidy
# Start server
go run main.go server
# Access swagger ui
http://localhost:8001/swagger
```

For more example code, please visit: https://github.com/jzero-io/examples

## Ecosystem

* jzero-intellij (jzero's goland plugin): https://github.com/jzero-io/jzero-intellij
* jzero-admin (Admin system based on jzero): https://github.com/jzero-io/jzero-admin
* templates (Jzero template market): https://templates.jzero.io

For more ecosystem, please visit: [https://docs.jzero.io/ecosystem/](https://docs.jzero.io/ecosystem/)

## Contributors

[Contribute](https://docs.jzero.io/community/contribute.html)

<a href="https://openomy.app/github/jzero-io/jzero" target="_blank" style="display: block; width: 100%;" align="center">
  <img src="https://openomy.app/svg?repo=jzero-io/jzero&chart=bubble&latestMonth=3" target="_blank" alt="Contribution Leaderboard" style="display: block; width: 100%;" />
</a>

## Stargazers over time

[![Star History Chart](https://api.star-history.com/svg?repos=jzero-io/jzero&type=Date)](https://star-history.com/#jzero-io/jzero&Date)

## Disclaimer

jzero is released under the MIT License and is provided completely free of charge. The authors and contributors assume no liability for any direct or indirect consequences arising from the use of this software, including but not limited to performance degradation, data loss, service interruptions, or any other type of damage.

No Warranty: This software comes with no express or implied warranties, including but not limited to fitness for a particular purpose, non-infringement, merchantability, and reliability.

User Responsibility: By using this software, you understand and agree to assume all risks and responsibilities associated with its use.
