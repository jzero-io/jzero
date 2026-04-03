---
home: false
icon: /icons/fluent-home-heart-20-filled.svg
title: Home
---

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/jzero.svg" style="width: 33%;" alt=""/>
</div>

## Introduction

**jzero** is a framework developed based on the [go-zero framework](https://github.com/zeromicro/go-zero) and [go-zero/goctl tools](https://github.com/zeromicro/go-zero/tree/master/tools/goctl). It enables one-click initialization of api/gateway/rpc projects.

Based on describable files (**api/proto/sql**), it automatically generates **server-side and client-side** framework code. Combined with built-in jzero-skills, it empowers AI to generate business logic code that follows best practices, reducing development cognitive load and freeing your hands!

Key features include:

* Supports flexible control of jzero configurations through combinations of **configuration files/command-line parameters/environment variables**, with minimal commands to generate code, AI-friendly
* Supports generating code based on **git-changed files** or specific descriptor files, or ignoring specific descriptor files, improving code generation efficiency for large projects
* Built-in common development templates with enhanced template features, supports **custom templates** to build proprietary enterprise code templates, significantly reducing development costs
* Supports **plugin architecture**, where functional modules can be dynamically loaded as independent plugins, supporting plugin creation, compilation, and uninstallation, perfectly adapted for team collaboration and module decoupling

For more details, please visit: [https://docs.jzero.io](https://docs.jzero.io)

## Design Philosophy

* **Developer Experience**: Provides a simple, easy-to-use, one-stop production-ready solution that enhances the development experience
* **Template-Driven**: All code generation is based on template rendering, with default generation following best practices, and supports custom template content
* **Ecosystem Compatibility**: Does not modify go-zero and go-zero/goctl, maintains ecosystem compatibility while addressing existing pain points and extending new features
* **Team Development**: Through module **layering** and **plugin** design, it's friendly to team development
* **Interface Design**: Does not depend on specific databases/cache/config centers and other infrastructure, allowing free choice based on actual needs

## Quick Start

::: code-tabs#shell
@tab jzero cli

```bash
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
# Visit swagger ui
http://localhost:8001/swagger
```

@tab jzero Docker

```bash
# One-click create project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project
# Download dependencies
go mod tidy
# Start server
go run main.go server
# Visit swagger ui
http://localhost:8001/swagger
```
:::
