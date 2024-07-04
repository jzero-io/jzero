# {{ .APP }}

## Install Jzero Framework

```shell
go install github.com/jzero-io/jzero@latest

jzero check
```

## Generate code

### Generate server code

```shell
jzero gen
```

### Generate swagger code

```shell
jzero gen swagger
```

you can see generated swagger json in `desc/swagger`

## Build docker image

```shell
# add a builder first
docker buildx create --use --name=mybuilder --driver docker-container --driver-opt image=dockerpracticesig/buildkit:master

# build and load
docker buildx build --platform linux/{{ .GoArch }} --progress=plain -t {{ .APP }}:latest . --load
```

## Documents

https://jzero.jaronnie.com