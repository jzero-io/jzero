APP: {{ .APP }}
Zrpc:
  ListenOn: 0.0.0.0:8000
  Mode: dev
  Name: {{ .APP }}.rpc
Gateway:
  Name: {{ .APP }}.gw
  Port: 8001
  Upstreams:
    - Grpc:
        Endpoints:
          - 0.0.0.0:8000
      Name: {{ .APP }}.gw
      ProtoSets:
        - .protosets/hello.pb

Log:
  encoding: plain
