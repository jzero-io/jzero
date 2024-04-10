## TODO: 
## 1. auto gen by dir proto
## 2. auto register rpc server
## 3. auto modify config.toml

# rpc
goctl rpc protoc jzerod/proto/credential.proto  -I./jzerod/proto --go_out=./jzerod --go-grpc_out=./jzerod  --zrpc_out=./jzerod --client=false -m --home .template
rm jzerod/credential.go

# rpc
goctl rpc protoc jzerod/proto/machine.proto  -I./jzerod/proto --go_out=./jzerod --go-grpc_out=./jzerod  --zrpc_out=./jzerod --client=false -m --home .template
rm jzerod/machine.go

# api
goctl api go --api jzerod/api/jzerod.api --dir ./jzerod --home .template

## rm etc
rm -rf jzerod/etc

# gen proto descriptor
protoc --include_imports -I./jzerod/proto --descriptor_set_out=protosets/credential.pb jzerod/proto/credential.proto
protoc --include_imports -I./jzerod/proto --descriptor_set_out=protosets/machine.pb jzerod/proto/machine.proto