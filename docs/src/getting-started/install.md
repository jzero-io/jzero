---
title: Install jzero
icon: /icons/marketeq-download-alt-4.svg
order: 2
---

## Install golang

Recommend using [gvm](https://github.com/jaronnie/gvm) to install golang environment

## Install jzero

Provides the following three ways to use jzero, choose one based on your actual situation

* Source installation(**go version >= go1.25.0**)
* Directly [download jzero binary](https://github.com/jzero-io/jzero/releases)
* Install jzero based on Docker, [image address](https://github.com/jzero-io/jzero/pkgs/container/jzero)

### Install jzero from source

```bash
# Set domestic proxy (optional)
# go env -w GOPROXY=https://goproxy.cn,direct
go install github.com/jzero-io/jzero/cmd/jzero@latest

# Get jzero version
jzero version

# Auto download required tools
jzero check
```

### Download jzero binary

[Download address](https://github.com/jzero-io/jzero/releases)

Select the corresponding package based on your operating system, extract and place in `$GOPATH/bin` directory

Execute the following to complete jzero environment setup

```shell
# Get jzero version
jzero version

# Auto download required tools
jzero check
```

### Install jzero based on Docker

```shell
# Get jzero version
docker run --rm ghcr.io/jzero-io/jzero:latest version
```

## Upgrade jzero

```shell
# Upgrade to latest version
jzero upgrade
# Upgrade to specific version
jzero upgrade --channel <commit_hash> or <tag>
```
