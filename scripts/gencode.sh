## TODO: 
## 1. auto gen by dir proto
## 2. auto register rpc server
## 3. auto modify worktab.toml

# rpc
goctl rpc protoc worktabd/proto/credential.proto  -I./worktabd/proto --go_out=./worktabd --go-grpc_out=./worktabd  --zrpc_out=./worktabd --client=false -m --home .template
rm worktabd/credential.go

# rpc
goctl rpc protoc worktabd/proto/machine.proto  -I./worktabd/proto --go_out=./worktabd --go-grpc_out=./worktabd  --zrpc_out=./worktabd --client=false -m --home .template
rm worktabd/machine.go

# api
goctl api go --api worktabd/api/worktabd.api --dir ./worktabd --home .template

## rm etc
rm -rf worktabd/etc

# gen proto descriptor
protoc --include_imports -I./worktabd/proto --descriptor_set_out=credential.pb worktabd/proto/credential.proto
protoc --include_imports -I./worktabd/proto --descriptor_set_out=machine.pb worktabd/proto/machine.proto