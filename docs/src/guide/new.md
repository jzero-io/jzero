---
title: 新建项目
icon: clone
order: 3
---

::: code-tabs#shell

@tab jzero

```bash
jzero new --module=github.com/jaronnie/app1 --dir=./app1 --app=app1
```

@tab Docker

```bash
docker run --rm \
  -v ./app1:/app/app1 jaronnie/jzero:latest \
  new --module=github.com/jaronnie/app1 \
  --dir=./app1 --app=app1
```
:::

flag 解释:

* module 表示新建项目的 go module
* dir 表示创建的项目目录路径
* app 表示项目名