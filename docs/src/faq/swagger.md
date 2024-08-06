---
title: swagger 问题
icon: vscode-icons:file-type-swagger
star: true
order: 1
category: faq
tag:
  - faq
---

## 1. 为什么我的 proto 字段是 order_id 但是生成的 swagger 是 orderId, 如何解决?

```protobuf
syntax = "proto3";

message GetOrderRequest {
    int32 order_id = 2;
}
```

生成的 pb.go 文件:

```go
type SayHelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	OrderId int32  `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}
```

swagger ui 如图所示:

<img src="https://oss.jaronnie.com/image-20240806110300854.png" style="zoom:67%;" />

解决方案:

```protobuf
syntax = "proto3";

message GetOrderRequest {
  int32 order_id = 2 [json_name = "order_id"];
}
```

<img src="https://oss.jaronnie.com/image-20240806110436742.png" alt="" style="zoom:67%;" />

如果需要自定义注释等信息, 可以使用

```protobuf
syntax = "proto3";
import "google/api/annotations.proto";
import "grpc-gateway/protoc-gen-openapiv2/options/annotations.proto";

message GetOrderRequest {
  int32 order_id = 2 [(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
    description: "订单 id"
  }, json_name = "order_id"];
}
```

<img src="https://oss.jaronnie.com/image-20240806110339702.png" style="zoom:67%;" />
