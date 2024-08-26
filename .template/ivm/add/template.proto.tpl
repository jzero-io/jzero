syntax = "proto3";

package {{ .Package }}{{ .Version }}pb;

import "third_party/google/api/annotations.proto";
import "third_party/validate/validate.proto";
import "third_party/grpc-gateway/protoc-gen-openapiv2/options/annotations.proto";

option go_package = "./pb/{{ .Package }}{{ .Version }}pb";

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    info: {
        version: "{{ .UrlVersion }}";
    };
};

{{range $v := .Methods}}message {{$v.Name | FirstUpper}}Request {}
message {{$v.Name | FirstUpper}}Response {}
{{end}}

{{range $service := .Services}}
service {{ $service | FirstUpper }}{{ $.Version | FirstUpper }} { {{range $m := $.Methods}}
    rpc {{ $m.Name }}({{$m.Name | FirstUpper}}Request) returns({{$m.Name | FirstUpper}}Response) {
        option (google.api.http) = {
            {{ $m.Verb }}: "/api/{{ $.UrlVersion }}/{{ $service }}/{{ $m.Name | FirstLower}}"
        };
    };
{{end}}
}
{{end}}
