{{ if .Serverless }}{{else}}syntax = "proto3";

package {{if .Serverless}}{{.APP | lower}}{{end}}version;
option go_package = "./types/{{if .Serverless}}{{.APP | lower}}{{end}}version";

message VersionRequest {}

message VersionResponse {
    string version = 1;
    string goVersion = 2;
    string commit = 3;
    string date = 4;
}

service {{if .Serverless}}{{.APP | ToCamel}}{{end}}Version {
    rpc Version(VersionRequest) returns(VersionResponse) {};
}{{end}}