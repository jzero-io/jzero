---
title: 新增可描述文件
icon: proicons:tag-add
order: 3.1
---

jzero 根据可描述文件(desc)生成代码:

* desc/api: api 可描述语言, 生成 http 服务端/客户端代码. [使用指南](../guide/api.md)
* desc/proto: proto 可描述语言, 生成 grpc 服务端/客户端代码. [使用指南](../guide/proto.md)
* desc/sql: sql 可描述语言, 生成数据库代码. [使用指南](../guide/model.md)

亦可基于配置生成 model 代码:

* model datasource: 通过远程数据源生成数据库代码. [使用指南](../guide/model.md)
* mongo type: 通过指定 mongo type 生成 mongodb 代码. [使用指南](../guide/mongo.md)

## 新增 api 文件

将在 desc/api 文件夹下新增 api 文件

```shell
# group 为 test
jzero add api test
# group 为 test/test1
jzero add api test/test1
```

## 新增 proto 文件

将在 desc/proto 文件夹下新增 proto 文件

```shell
# Service 为 Test
jzero add proto test
# Service 为 TestTest1
jzero add proto test/test1
```

## 新增 sql 文件

将在 desc/sql 文件夹下新增 sql 文件

```shell
# table 名为 test
jzero add sql test
```
