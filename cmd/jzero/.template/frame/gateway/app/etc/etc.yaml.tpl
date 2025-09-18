zrpc:
  listenOn: 0.0.0.0:8000
  mode: dev
  name: {{ .APP }}.rpc
gateway:
  name: {{ .APP }}.gw
  port: 8001
  upstreams:
    - grpc:
        endpoints:
          - 0.0.0.0:8000
      name: {{ .APP }}.gw

log:
    serviceName: {{ .APP }}
    encoding: plain
    level: info
    mode: console
{{ if has "model" .Features }}
sqlx:
    driverName: "mysql"
    dataSource: "root:123456@tcp(127.0.0.1:3306)/{{ .APP }}?charset=utf8mb4&parseTime=True&loc=Local"
{{ end }}{{ if has "redis" .Features }}
redis:
    host: "127.0.0.1:6379"
    type: "node"
    pass: "123456"{{ end }}