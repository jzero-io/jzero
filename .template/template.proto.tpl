syntax = "proto3";

package {{ .Service }}{{ .Version }}pb;

import "google/api/annotations.proto";
import "validate/validate.proto";

option go_package = "./pb/{{ .Service }}{{ .Version }}pb";

{{range $v := .Methods}}message {{$v.Name | FirstUpper}}Request {}
message {{$v.Name | FirstUpper}}Response {}
{{end}}

service {{ .Service }}{{ .Version }} { {{range $v := .Methods}}
    rpc {{ $v.Name }}({{$v.Name | FirstUpper}}Request) returns({{$v.Name | FirstUpper}}Response) {
        option (google.api.http) = {
            {{ $v.Verb }}: "/api/{{ $.UrlVersion }}/{{ $.Service }}/{{ $v.Name | FirstLower}}"
        };
    };
{{end}}
}
