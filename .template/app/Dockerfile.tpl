FROM --platform=$BUILDPLATFORM golang:alpine as builder

ARG TARGETARCH
ARG LDFLAGS

ENV GOPROXY https://goproxy.cn,direct

WORKDIR /usr/local/go/src/app

COPY ./ ./

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -a -ldflags="$LDFLAGS" -o /dist/app main.go \
    && cp -r etc /dist/etc \
    && find desc/proto -type f -name '*.pb' | while read file; do mkdir -p /dist/$(dirname $file) && cp $file /dist/$file; done


FROM --platform=$TARGETPLATFORM alpine:latest

WORKDIR /dist

COPY --from=builder /dist .

EXPOSE 8000 8001

CMD ["./app", "server"]