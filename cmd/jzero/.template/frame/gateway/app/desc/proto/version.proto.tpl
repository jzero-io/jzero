{{ if .Serverless }}{{else}}syntax = "proto3";

package version;

import "google/api/annotations.proto";
import "grpc-gateway/protoc-gen-openapiv2/options/annotations.proto";

option go_package = "./types/version";

message VersionRequest {}

message VersionResponse {
    string version = 1;
    string goVersion = 2;
    string commit = 3;
    string date = 4;
}

service Version {
    rpc Version(VersionRequest) returns(VersionResponse) {
        option (google.api.http) = {
            get: "/api/version"
        };
    };
}{{end}}