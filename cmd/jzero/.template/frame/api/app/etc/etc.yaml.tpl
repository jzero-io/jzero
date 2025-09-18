rest:
    name: {{ .APP }}-api
    host: 0.0.0.0
    port: 8001

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