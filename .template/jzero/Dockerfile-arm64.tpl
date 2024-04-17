FROM arm64v8/alpine:latest

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

RUN apk add tzdata ca-certificates curl bash
ENV TZ Asia/Shanghai

WORKDIR /app
COPY dist/{{ .APP }}_linux_arm64/{{ .APP }} /app/{{ .APP }}
COPY config.toml /app/config.toml
COPY .protosets /app/.protosets

EXPOSE 8001 8002

ENTRYPOINT ["./{{ .APP }}", "daemon", "--config", "config.toml"]