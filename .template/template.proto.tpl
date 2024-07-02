syntax = "proto3";

package {{ .Package }}{{ .Version }}pb;

import "google/api/annotations.proto";
import "validate/validate.proto";

option go_package = "./pb/{{ .Package }}{{ .Version }}pb";

{{range $v := .Methods}}message {{$v.Name | FirstUpper}}Request {}
message {{$v.Name | FirstUpper}}Response {}
{{end}}

{{range $service := .Services}}
service {{ $service }}{{ $.Version }} { {{range $m := $.Methods}}
    rpc {{ $m.Name }}({{$m.Name | FirstUpper}}Request) returns({{$m.Name | FirstUpper}}Response) {
        option (google.api.http) = {
            {{ $m.Verb }}: "/api/{{ $.UrlVersion }}/{{ $service }}/{{ $m.Name | FirstLower}}"
        };
    };
{{end}}
}
{{end}}
