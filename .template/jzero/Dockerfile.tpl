FROM alpine:latest

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

RUN apk add tzdata ca-certificates curl bash
ENV TZ Asia/Shanghai

WORKDIR /app
COPY dist/{{ .APP }}_linux_amd64_v1/{{ .APP }} /app/{{ .APP }}
COPY config.{{ .ConfigType }} /app/config.{{ .ConfigType }}
COPY .protosets /app/.protosets

EXPOSE 8001 8002

ENTRYPOINT ["./{{ .APP }}", "daemon"]