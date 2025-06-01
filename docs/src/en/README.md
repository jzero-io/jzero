---
home: false
icon: fluent:home-heart-20-filled
title: Home  
---

<div style="text-align: center;">  
  <img src="https://oss.jaronnie.com/jzero.jpg" style="width: 33%;" alt=""/>  
</div>  

## Introduction

`jzero` is a low-code microservice development framework based on [go-zero](https://go-zero.dev). It automatically generates **server-side code/client-side code/database** code through descriptive files (**api/proto/sql**), reducing development complexity and freeing your hands!

jzero has the following features:

* Supports controlling command parameters through a combination of configuration files, command-line arguments, and environment variables, eliminating tedious command configurations
* Supports generating code incrementally based on git-tracked file changes, significantly improving code generation efficiency for large projects
* Optimizes existing pain points in go-zero and extends new features
* Built-in commonly used development templates with enhanced template capabilities, supporting custom template content to build enterprise-level code templates
* All supporting tools are cross-platform, compatible with Windows/Mac/Linux

## Quick Start

:::tip Windows users, please execute all commands in PowerShell  
:::

::: code-tabs#shell  
@tab jzero

```bash  
# Install jzero  
go install github.com/jzero-io/jzero/cmd/jzero@latest  
# One-click installation of required tools  
jzero check
# Create a project
jzero new your_project
cd your_project
# Download dependencies  
go mod tidy  
# Generate swagger  
jzero gen swagger  
# Start the server
go run main.go server  
# Access swagger UI  
http://localhost:8001/swagger  
```

@tab Docker

```bash
# Create a project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project  
cd your_project
# Download dependencies
go mod tidy
# Generate swagger
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# Start the server
go run main.go server
# Access swagger UI
http://localhost:8001/swagger  
```
:::