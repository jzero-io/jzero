---
title: model 数据库
icon: puzzle-piece
star: true
order: 5
category: 开发
tag:
  - Guide
---

jzero 推荐使用 go-zero sqlx 完成对数据库的 crud 操作.

jzero 数据库规范:

* sql 文件放在 daemon/desc/sql
* 生成的 model 放在 daemon/model

jzero gen 时会自动检测 daemon/desc/sql 下的 sql 文件, 并将生成的 model 放在 daemon/model 下