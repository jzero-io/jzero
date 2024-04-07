# worktab

work table

## 技术栈

* cobra(实现 worktab cli)
* go-zero(实现 worktabd)
* grpc-gateway(提供 grpc 的 http 接口)
* gin
* protoc-gen-go-httpsdk(自动生成 worktabd 的 go sdk 库)
* vue(worktab 内置 vue 实现的 worktab ui)

## worktab worktabd

worktabd 即 worktab 的服务端, 支持的协议:
* unix-socket
* tcp
* ssh

```shell
worktab worktabd
```
