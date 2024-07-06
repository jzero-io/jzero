---
title: 生成 swagger 文档
icon: lightbulb
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
docker run --rm -v ${PWD}:/app jaronnie/jzero:latest gen swagger
```
:::

## 在线访问 swagger ui

ip:port/swagger