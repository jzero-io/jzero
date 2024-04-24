APP = "{{ .APP }}"
Name = "{{ .APP }}.rpc"
ListenOn = "0.0.0.0:8000"
Key = "{{ .APP }}.rpc"
Mode = "dev"

[Log]
  ServiceName = "{{ .APP }}"
  Level = "info"
  Mode = "file"
  encoding = "plain"
  KeepDays = 30
  MaxBackups = 7
  MaxSize = 50
  Rotation = "size"

[DevServer]
Enabled = true

[Gateway]
Name = "{{ .APP }}.gw"
Port = 8001

  [[Gateway.Upstreams]]
  Name = "{{ .APP }}.gw"
  ProtoSets = [ ".protosets/credential.pb" ]

    [Gateway.Upstreams.Grpc]
    Endpoints = [ "0.0.0.0:8000" ]