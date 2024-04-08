# rpc
goctl rpc protoc worktabd/proto/worktabd.proto  -I./worktabd/proto --go_out=./worktabd --go-grpc_out=./worktabd  --zrpc_out=./worktabd --client=false

# api
goctl api go --api worktabd/proto/worktabd.api --dir ./worktabd

# gen proto descriptor
protoc --include_imports -I./worktabd/proto --descriptor_set_out=worktabd.pb worktabd/proto/worktabd.proto