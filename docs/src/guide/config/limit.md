---
title: 限流配置
icon: gears
star: true
order: 5
category: 配置
tag:
  - Guide
---

修改 etc/etc.yaml, 增加以下配置, 设置最大 qps 100

```yaml
Zrpc:
  MaxConns: 100
Gateway:
  MaxConns: 100
```

由于 jzero 默认集成了 go-zero 三个特性

* rpc
* api
* gateway

我们依次测试这三种类型的接口

::: tip
https://github.com/zeromicro/go-zero/issues/4097

两种路由的限流都没生效

* api 生成的 handler 注册到 gateway server 后
* gateway server AddRoute 的路由
:::

```shell
# test rpc
ghz --insecure -c 110 -n 110  \
  --call credentialpb.credential.CredentialVersion \
  0.0.0.0:8000
  
$ ghz --insecure -c 110 -n 110  \
  --call credentialpb.credential.CredentialVersion \
  0.0.0.0:8000

Summary:
  Count:	110
  Total:	117.02 ms
  Slowest:	106.65 ms
  Fastest:	51.92 ms
  Average:	77.79 ms
  Requests/sec:	940.03

Response time histogram:
  51.918  [1]  |∎∎
  57.391  [4]  |∎∎∎∎∎∎∎∎∎
  62.865  [3]  |∎∎∎∎∎∎∎
  68.339  [8]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  73.813  [9]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  79.286  [6]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  84.760  [11] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  90.234  [12] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  95.707  [17] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  101.181 [16] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  106.655 [13] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎

Latency distribution:
  10 % in 65.12 ms
  25 % in 73.70 ms
  50 % in 88.74 ms
  75 % in 96.70 ms
  90 % in 102.15 ms
  95 % in 104.78 ms
  99 % in 105.13 ms

Status code distribution:
  [Unavailable]   10 responses
  [OK]            100 responses

Error distribution:
  [10]   rpc error: code = Unavailable desc = concurrent connections over limit

# test api
hey -z 1s -c 120 -q 1 'http://localhost:8001/api/v1/hello/you'

Summary:
  Total:	1.0821 secs
  Slowest:	0.0745 secs
  Fastest:	0.0196 secs
  Average:	0.0475 secs
  Requests/sec:	110.8997

  Total data:	8880 bytes
  Size/request:	74 bytes

Response time histogram:
  0.020 [1]	|■■
  0.025 [4]	|■■■■■■■■
  0.031 [11]	|■■■■■■■■■■■■■■■■■■■■■■
  0.036 [10]	|■■■■■■■■■■■■■■■■■■■■
  0.042 [15]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.047 [15]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.053 [19]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.058 [13]	|■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.064 [20]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.069 [8]	|■■■■■■■■■■■■■■■■
  0.075 [4]	|■■■■■■■■


Latency distribution:
  10% in 0.0297 secs
  25% in 0.0379 secs
  50% in 0.0487 secs
  75% in 0.0584 secs
  90% in 0.0635 secs
  95% in 0.0670 secs
  99% in 0.0745 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0143 secs, 0.0196 secs, 0.0745 secs
  DNS-lookup:	0.0072 secs, 0.0024 secs, 0.0108 secs
  req write:	0.0005 secs, 0.0000 secs, 0.0031 secs
  resp wait:	0.0250 secs, 0.0026 secs, 0.0457 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0010 secs

Status code distribution:
  [200]	120 responses


# test gateway
# 用 hey 工具来进行压测，压测 120 个并发，执行 1 秒, 有 20 个被限流
hey -z 1s -c 120 -q 1 'http://localhost:8001/api/v1/credential/version'

Summary:
  Total:	1.1574 secs
  Slowest:	0.1511 secs
  Fastest:	0.0217 secs
  Average:	0.1111 secs
  Requests/sec:	103.6849

  Total data:	5800 bytes
  Size/request:	48 bytes

Response time histogram:
  0.022 [1]	|■
  0.035 [17]	|■■■■■■■■■■■■■■■
  0.048 [2]	|■■
  0.061 [0]	|
  0.073 [2]	|■■
  0.086 [2]	|■■
  0.099 [0]	|
  0.112 [5]	|■■■■■
  0.125 [24]	|■■■■■■■■■■■■■■■■■■■■■■
  0.138 [44]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.151 [23]	|■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  10% in 0.0330 secs
  25% in 0.1144 secs
  50% in 0.1264 secs
  75% in 0.1366 secs
  90% in 0.1424 secs
  95% in 0.1438 secs
  99% in 0.1511 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0128 secs, 0.0217 secs, 0.1511 secs
  DNS-lookup:	0.0046 secs, 0.0009 secs, 0.0079 secs
  req write:	0.0004 secs, 0.0000 secs, 0.0023 secs
  resp wait:	0.0969 secs, 0.0056 secs, 0.1301 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0002 secs

Status code distribution:
  [200]	100 responses
  [503]	20 responses
```

