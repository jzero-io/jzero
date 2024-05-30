FROM --platform=$BUILDPLATFORM golang:1.22 as builder

ARG TARGETARCH
ARG LDFLAGS

WORKDIR /usr/local/go/src/app
COPY ./ ./
ENV GOPROXY https://goproxy.io,direct
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -a -ldflags="$LDFLAGS" -o /app main.go


FROM --platform=$TARGETPLATFORM alpine:latest
WORKDIR /app
COPY --from=builder /app .
COPY etc/etc.yaml /app/etc/etc.yaml
COPY .protosets /app/.protosets
EXPOSE 8000 8001
ENTRYPOINT ["./app"]
CMD ["server"]