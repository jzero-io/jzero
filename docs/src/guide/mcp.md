---
title: mcp server
icon: tabler:photo-ai
order: 5.3
category: 开发
tag:
  - Guide
---

## 使用 mcp

使用 mcp client, 如 [deepchat](https://deepchat.thinkinai.xyz/)

![](http://oss.jaronnie.com/image-20250512113200546.png)

## 测试 jzero mcp 之底层协议

```shell
$ jzero mcp test
[INPUT] Enter your command (press Enter to send):
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"roots":{"listChanged":true},"sampling":{}},"clientInfo":{"name":"ExampleClient","version":"1.0.0"}}}
[SERVER OUTPUT] {"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"logging":{},"resources":{"subscribe":true,"listChanged":true},"tools":{}},"serverInfo":{"name":"Used to create project by templates and generate server/client code by proto and api file.\n","version":"1.0.0"}}}
[INPUT] Enter your command (press Enter to send):
{"jsonrpc":"2.0","method":"notifications/initialized"}
[INPUT] Enter your command (press Enter to send):
{"jsonrpc":"2.0","id":2,"method":"tools/list"}
```

