APP = "{{ .APP }}"
Name = "{{ .APP }}.rpc"
ListenOn = "0.0.0.0:8000"
Key = "{{ .APP }}.rpc"
Mode = "dev"

[DevServer]
Enabled = true

[Gateway]
Name = "{{ .APP }}.gw"
Port = 8001

  [[Gateway.Upstreams]]
  Name = "{{ .APP }}.gw"
  ProtoSets = [ ".protosets/credential.pb", ".protosets/machine.pb" ]

    [Gateway.Upstreams.Grpc]
    Endpoints = [ "0.0.0.0:8000" ]