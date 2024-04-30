---
title: go mod tidy error when using docker
icon: bug
order: 1
---

![](https://oss.jaronnie.com/image-20240430144344653.png)

docker 采用的是最新的 Go 版本, 而本地使用的是较旧的版本

两种解决办法:

* 修改 go.mod 文件, 将 go 1.22.2 改为本地的 go 版本, 如 go 1.19
* 升级本地 go 版本为最新版本