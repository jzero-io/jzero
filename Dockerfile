FROM golang:alpine

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

RUN apk update --no-cache \
    && apk add --no-cache tzdata ca-certificates curl bash protoc

RUN go install github.com/zeromicro/go-zero/tools/goctl@latest \
  && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV TZ Asia/Shanghai

WORKDIR /app
COPY dist/jzero_linux_amd64_v1/jzero /app/jzero
COPY config.toml /app/config.toml
COPY .protosets /app/.protosets

EXPOSE 8000 8001

ENTRYPOINT ["./jzero"]