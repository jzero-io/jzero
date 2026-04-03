---
title: Add descriptor files
icon: /icons/proicons-tag-add.svg
order: 3.1
---

jzero generates code based on descriptor files (desc):

* desc/api: api descriptor language, generate http server/client code. [User guide](../guide/api.md)
* desc/proto: proto descriptor language, generate grpc server/client code. [User guide](../guide/proto.md)
* desc/sql: sql descriptor language, generate database code. [User guide](../guide/model.md)

Can also generate model code based on configuration:

* model datasource: Generate database code through remote datasource. [User guide](../guide/model.md)
* mongo type: Generate mongodb code by specifying mongo type. [User guide](../guide/mongo.md)

## Add api file

Will add api file under desc/api folder

```shell
# group is test
jzero add api test
# group is test/test1
jzero add api test/test1
```

## Add proto file

Will add proto file under desc/proto folder

```shell
# Service is Test
jzero add proto test
# Service is TestTest1
jzero add proto test/test1
```

## Add sql file

Will add sql file under desc/sql folder

```shell
# table name is test
jzero add sql test
```
