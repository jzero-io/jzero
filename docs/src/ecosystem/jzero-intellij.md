---
title: jzero-intellij(goland 插件)
icon: catppuccin:folder-admin
order: 2
---

goland + jzero 插件, 能极大提升开发体验感

## 安装 jzero 插件

[下载地址](https://github.com/jzero-io/jzero-intellij/releases)

![](https://oss.jaronnie.com/image-20251217111538091.png)

## 功能特性

* 新增 api/proto/sql 文件
* api 文件高亮/跳转等
* api/proto 文件跳转至 logic 文件
* logic 文件跳转至 api/proto 文件
* api/proto/sql 文件行首增加执行按钮生成代码
* .jzero.yaml 文件增加执行按钮生成代码

## 效果展示

<video width="720" height="450" controls>
  <source src="https://oss.jaronnie.com/jzero-intellij.mp4" type="video/mp4">
</video>

## release

### v1.0.0(2026.01.01)

### v1.0.1(2026.01.05)

* 优化 api 文件跳转到 logic

### v1.1.0

:::important requires jzero >= v1.1.0
:::

* 支持从 logic 文件跳转到 api/proto 文件