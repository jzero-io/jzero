FROM golang:alpine

ENV CGO_ENABLED=0

LABEL \
  org.opencontainers.image.title="jzero" \
  org.opencontainers.image.description="jzero framework" \
  org.opencontainers.image.url="https://github.com/jzero-io/jzero" \
  org.opencontainers.image.documentation="https://github.com/jzero-io/jzero#readme" \
  org.opencontainers.image.source="https://github.com/jzero-io/jzero" \
  org.opencontainers.image.licenses="MIT" \
  maintainer="jaronnie <jaron@jaronnie.com>"

WORKDIR /app

COPY dist/jzero_linux_amd64_v1/jzero /dist/jzero_linux_amd64/jzero
COPY dist/jzero_linux_arm64_v8.0/jzero /dist/jzero_linux_arm64/jzero

RUN if [ `go env GOARCH` = "amd64" ]; then \
      cp /dist/jzero_linux_amd64/jzero /usr/local/bin/jzero; \
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