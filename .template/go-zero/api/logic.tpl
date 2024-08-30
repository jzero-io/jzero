package {{.pkgName}}

import (
        {{.imports}}

        {{if or (eq .request "") (eq .responseType "error")}}"net/http"{{end}}
)

type {{.logic}} struct {
        logx.Logger
        ctx    context.Context
        svcCtx *svc.ServiceContext
        {{if eq .request ""}}r *http.Request{{end}}
        {{if eq .responseType "error"}}w http.ResponseWriter{{end}}
}

{{if .hasDoc}}{{.doc}}{{end}}
func New{{.logic}}(ctx context.Context, svcCtx *svc.ServiceContext, {{if eq .request ""}}r *http.Request, {{end}}{{if eq .responseType "error"}}w http.ResponseWriter{{end}}) *{{.logic}} {
        return &{{.logic}}{
                Logger: logx.WithContext(ctx),
                ctx:    ctx,
                svcCtx: svcCtx,
                {{if eq .request ""}}r: r,{{end}}
                {{if eq .responseType "error"}}w: w,{{end}}
        }
}

func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
        // todo: add your logic here and delete this line

        {{.returnString}}
}