APP = "{{ .APP }}"
Name = "{{ .APP }}.rpc"
ListenOn = "0.0.0.0:8000"
Key = "{{ .APP }}.rpc"
Mode = "dev"

[DevServer]
Enabled = true

# Grpc 链路追踪
# [Telemetry]
# Name = "{{ .APP }}"
# Endpoint = "http://127.0.0.1:14268/api/traces"
# Batcher = "jaeger"
# Sampler = 1.0

[Gateway]
Name = "{{ .APP }}.gw"
Port = 8001

# Gateway 链路追踪
# [Gateway.Telemetry]
# Name = "{{ .APP }}"
# Endpoint = "http://127.0.0.1:14268/api/traces"
# Batcher = "jaeger"
# Sampler = 1.0

  [[Gateway.Upstreams]]
  Name = "{{ .APP }}.gw"
  ProtoSets = [ ".protosets/credential.pb", ".protosets/machine.pb" ]

    [Gateway.Upstreams.Grpc]
    Endpoints = [ "0.0.0.0:8000" ]