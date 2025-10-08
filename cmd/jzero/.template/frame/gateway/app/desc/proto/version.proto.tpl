{{ if .Serverless }}{{else}}syntax = "proto3";

package versionpb;

import "google/api/annotations.proto";
import "grpc-gateway/protoc-gen-openapiv2/options/annotations.proto";

option go_package = "./pb/versionpb";

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
            get: "/version"
        };
    };
}{{end}}