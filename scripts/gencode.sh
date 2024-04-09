# rpc
goctl rpc protoc worktabd/proto/credential.proto  -I./worktabd/proto --go_out=./worktabd --go-grpc_out=./worktabd  --zrpc_out=./worktabd --client=false -m 

# rpc
goctl rpc protoc worktabd/proto/machine.proto  -I./worktabd/proto --go_out=./worktabd --go-grpc_out=./worktabd  --zrpc_out=./worktabd --client=false -m

# api
goctl api go --api worktabd/proto/worktabd.api --dir ./worktabd

# gen proto descriptor
protoc --include_imports -I./worktabd/proto --descriptor_set_out=credential.pb worktabd/proto/credential.proto
protoc --include_imports -I./worktabd/proto --descriptor_set_out=machine.pb worktabd/proto/machine.proto