// @ts-ignore
import { RequestMethod } from "umi-request";

{{range $scope := .Scopes}}import { {{$scope | FirstUpper }}Client } from "./typed/{{$scope}}/{{$scope}}_client";
{{end}}

export default class Clientset {
    constructor(private request: RequestMethod) { }
{{range $scope := .Scopes}}
	get {{$scope | FirstUpper}}(): {{$scope | FirstUpper}}Client {
		return new {{$scope | FirstUpper}}Client(this.request);
	}
{{end}}
}