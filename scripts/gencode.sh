## TODO: 
## 1. auto gen by dir proto
## 2. auto register rpc server
## 3. auto modify config.toml

# rpc
goctl rpc protoc daemon/proto/credential.proto  -I./daemon/proto --go_out=./daemon --go-grpc_out=./daemon  --zrpc_out=./daemon --client=false -m --home .template
rm daemon/credential.go

# rpc
goctl rpc protoc daemon/proto/machine.proto  -I./daemon/proto --go_out=./daemon --go-grpc_out=./daemon  --zrpc_out=./daemon --client=false -m --home .template
rm daemon/machine.go

# api
goctl api go --api daemon/api/jzero.api --dir ./daemon --home .template

## rm etc
rm -rf daemon/etc
rm daemon/jzero.go

# gen proto descriptor
protoc --include_imports -I./daemon/proto --descriptor_set_out=.protosets/credential.pb daemon/proto/credential.proto
protoc --include_imports -I./daemon/proto --descriptor_set_out=.protosets/machine.pb daemon/proto/machine.proto