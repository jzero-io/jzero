# {{ .APP }}

## Install Jzero Framework

```shell
go install github.com/jzero-io/jzero/cmd/jzero@latest

jzero check
```

## Generate code

### Generate server code

```shell
jzero gen
```

## Build docker image

```shell
# add a builder first
docker buildx create --use --name=mybuilder --driver docker-container --driver-opt image=dockerpracticesig/buildkit:master

# build and load
docker buildx build --platform linux/{{ .GoArch }} --progress=plain -t {{ .APP }}:latest . --load
```

## Documents

https://docs.jzero.io