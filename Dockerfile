FROM --platform=$TARGETPLATFORM golang:alpine

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct
WORKDIR /app
COPY config.toml /app/config.toml

RUN if [ `go env GOARCH` = "amd64" ]; then \
      cp /dist/jzero_linux_amd64_v1/jzero /app/jzero; \
    elif [ `go env GOARCH` = "arm64" ]; then \
      cp /dist/jzero_linux_arm64/jzero /app/jzero; \
    fi

RUN apk update --no-cache \
  && apk add --no-cache tzdata ca-certificates protoc

RUN /app/jzero check

RUN rm -rf /dist \
    && rm -rf /go/pkg/mod \
    && rm -rf /go/pkg/sumdb

EXPOSE 8000 8001
ENTRYPOINT ["./jzero"]
CMD ["-h"]