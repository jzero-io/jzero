rest:
    name: {{ .APP }}-api
    host: 0.0.0.0
    port: 8001

log:
    serviceName: {{ .APP }}
    encoding: plain
    level: info
    mode: console