{{- if not .EnableStylingCheck}}
/* eslint-disable */
// @ts-nocheck
{{- end}}
/*
* This file is a generated Typescript file for SDK Gateway, DO NOT MODIFY
*/

import { RequestMethod } from "umi-request";
import { Request, Response } from "../../rest/request";

/*
    TODO: add other import
*/

export class {{.Resource | FirstUpper}}Client {
    constructor(private request: RequestMethod) { }
{{range $interface := .HTTPInterfaces}}
{{end}}
}