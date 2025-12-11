syntax = "proto3";

package {{ .Package | base }};

option go_package = "./types/{{ .Package }}";

import "google/api/annotations.proto";
import "grpc-gateway/protoc-gen-openapiv2/options/annotations.proto";

message CreateRequest {}

message CreateResponse {}

service {{ .Service | FirstUpper }} {
    rpc Create(CreateRequest) returns(CreateResponse) {
        option (google.api.http) = {
            post: "/api/{{ .Package }}"
            body: "*"
        };
    };
}