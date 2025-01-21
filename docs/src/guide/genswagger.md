---
title: 生成 swagger 文档
icon: vscode-icons:file-type-swagger
order: 5.1
---

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen swagger
```

@tab Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
```
:::

## 在线访问 swagger ui

swagger ui 地址: **localhost:8001/swagger**

![](https://oss.jaronnie.com/image-20240731134511973.png)