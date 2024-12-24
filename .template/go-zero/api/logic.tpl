package {{.pkgName}}

import (
        "net/http"

        {{.imports}}
)

type {{.logic}} struct {
        logx.Logger
        ctx    context.Context
        svcCtx *svc.ServiceContext
        r *http.Request
        {{if eq .responseType "error"}}w http.ResponseWriter{{end}}
}

{{if .hasDoc}}{{.doc}}{{end}}
func New{{.logic}}(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request, {{if eq .responseType "error"}}w http.ResponseWriter{{end}}) *{{.logic}} {
        return &{{.logic}}{
                Logger: logx.WithContext(ctx),
                ctx:    ctx,
                svcCtx: svcCtx,
                r: r,
                {{if eq .responseType "error"}}w: w,{{end}}
        }
}

func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
        // todo: add your logic here and delete this line

        {{.returnString}}
}