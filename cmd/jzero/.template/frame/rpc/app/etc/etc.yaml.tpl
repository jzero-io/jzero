zrpc:
    listenOn: 0.0.0.0:8000
    mode: dev
    name: {{ .APP }}.rpc

log:
    serviceName: {{ .APP }}
    encoding: plain
    level: info
    mode: console