syntax = "proto3";

package {{ .Package | base }};

option go_package = "./types/{{ .Package }}";

message CreateRequest {}

message CreateResponse {}

service {{ .Service | FirstUpper }} {
    rpc Create(CreateRequest) returns(CreateResponse) {};
}