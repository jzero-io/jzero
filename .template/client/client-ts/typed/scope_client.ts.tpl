import { RequestMethod } from "umi-request";
{{range $resource := .Resources}}import { {{$resource | FirstUpper}}Client } from "./{{$resource}}";
{{end}}

export class {{.Scope | FirstUpper}}Client {
	constructor(private request: RequestMethod) { }

{{range $resource := .Resources}}
	get {{$resource | FirstUpper}}(): {{$resource | FirstUpper}}Client {
		return new {{$resource | FirstUpper}}Client(this.request);
	}
{{end}}
}