FROM --platform=$TARGETPLATFORM golang:alpine

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

WORKDIR /app

COPY config.toml /root/.jzero/config.toml
COPY dist/jzero_linux_amd64_v1/jzero /dist/jzero_linux_amd64_v1/jzero
COPY dist/jzero_linux_arm64/jzero /dist/jzero_linux_arm64/jzero

RUN if [ `go env GOARCH` = "amd64" ]; then \
      cp /dist/jzero_linux_amd64_v1/jzero /usr/local/bin/jzero; \
    elif [ `go env GOARCH` = "arm64" ]; then \
      cp /dist/jzero_linux_arm64/jzero /usr/local/bin/jzero; \
    fi

RUN apk update --no-cache \
  && apk add --no-cache tzdata ca-certificates protoc \
  && jzero check \
  && rm -rf /dist \
  && rm -rf /go/pkg/mod \
  && rm -rf /go/pkg/sumdb

ENTRYPOINT ["jzero"]