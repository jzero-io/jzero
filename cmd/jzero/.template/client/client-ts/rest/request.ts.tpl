import { RequestMethod } from "umi-request"

type PathParam<T> = {
    name: string
    value: T
}

type QueryParam = {
    name: string
    value: string
}

export type Response<T> = {
    code: number,
    message: string,
    data:  T
}


type Verb = "POST" | "GET"

export class Request<R, D=any> {
    request: RequestMethod
    verb: Verb
    subPath: string
    params: string
    body: D
    headers: Record<string, string>

    constructor(request: RequestMethod) {
        this.request = request;
    }

    setVerb(verb: Verb): this {
        this.verb = verb
        return this
    }

    setSubPath<T>(subPath: string, ...args: PathParam<T>[]): this {
        for (let v of args) {
            let placeHolder = "{"+v.name+"}";
            subPath = subPath.replace(placeHolder, String(v.value))
        }
        this.subPath = subPath
        return this
    }

    setHeaders(headers: Record<string, string>): this {
        this.headers = headers
        return this
    }

    setParams(...args: QueryParam[]): this {
        if (args.length == 0) {
            return this
        }

        let queryParams: string[] = []
        for (let arg of args) {
            let queryParam = arg.name + "=" + arg.value;
            queryParams.push(queryParam)
        }
        this.params = queryParams.join("&")
        return this
    }

    setBody(body: D): this {
        this.body = body
        return this
    }

    // @ts-ignore
    async send(): Promise<Response<R>> {
        if (this.verb == null) {
            throw new Error("verb is not set")
        }

        let url = this.subPath;

        return this.request<Response<R>>(url, {
            method: this.verb,
            data: this.body,
            headers: this.headers
        }).then(resp => {
            // @ts-ignore
            return Promise.resolve(resp)
        })
    }
}