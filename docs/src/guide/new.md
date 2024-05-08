---
title: 新建项目
icon: clone
order: 3
---

::: code-tabs#shell

@tab jzero

```bash
jzero new app1
```

@tab Docker(amd64)

```bash
docker run --rm -v ${PWD}/app1:/app/app1 jaronnie/jzero:latest new app1
```

@tab Docker(arm64)

```bash
docker run --rm -v ${PWD}/app1:/app/app1 jaronnie/jzero:latest-arm64 new app1
```
:::

all flags:

* module 设置 go module
* dir 设置生成的项目路径