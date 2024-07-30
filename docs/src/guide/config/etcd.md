---
title: etcd 配置
icon: carbon:database-etcd
star: true
order: 7
category: 配置
tag:
  - Guide
---

## zrpc

修改 etc/etc.yaml, 增加以下配置, 设置 etcd endpoints 和该服务在 etcd 上的 key

```yaml
zrpc:
    etcd:
        key: your_project.rpc
        hosts:
            - 127.0.0.1:2379
    listenOn: 0.0.0.0:8000
    mode: dev
    name: your_project.rpc
```

## 快速部署测试使用的 etcd 环境

```shell
docker pull bitnami/etcd:3.5.14
# 如果无法 pull
docker pull registry.cn-hangzhou.aliyuncs.com/jaronnie/etcd:3.5.14
docker tag registry.cn-hangzhou.aliyuncs.com/jaronnie/etcd:3.5.14 bitnami/etcd:3.5.14
```

```shell
docker run -p 2379:2379 -e ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd:3.5.14 
```

启动服务端后，查看 etcd 上的注册信息, 以下表示成功

![](https://oss.jaronnie.com/image-20240710222837633.png)

