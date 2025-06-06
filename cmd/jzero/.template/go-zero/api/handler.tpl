// Code generated by goctl. Templates Edited by jzero. DO NOT EDIT.

package {{.PkgName}}

import (
        "net/http"

        "github.com/zeromicro/go-zero/rest/httpx"

        {{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
        {{ if and .HasRequest .HasResp }}return func(w http.ResponseWriter, r *http.Request) {
                var req types.{{.RequestType}}
                if err := httpx.Parse(r, &req); err != nil {
                        httpx.ErrorCtx(r.Context(), w, err)
                        return
                }

                l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx, r)
                resp, err := l.{{.Call}}(&req)
                if err != nil {
                        httpx.ErrorCtx(r.Context(), w, err)
                } else {
                        httpx.OkJsonCtx(r.Context(), w, resp)
                }
        } {{else if and .HasRequest (not .HasResp)}}return func(w http.ResponseWriter, r *http.Request) {
                var req types.{{.RequestType}}
                if err := httpx.Parse(r, &req); err != nil {
                        httpx.ErrorCtx(r.Context(), w, err)
                        return
                }

                l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx, r, w)
                err := l.{{.Call}}(&req)
                if err != nil {
                        httpx.ErrorCtx(r.Context(), w, err)
                } else {
                        httpx.Ok(w)
                }
        } {{else if and (not .HasRequest) .HasResp}}return func(w http.ResponseWriter, r *http.Request) {
                l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx, r)
                resp, err := l.{{.Call}}()
                if err != nil {
                        httpx.ErrorCtx(r.Context(), w, err)
                } else {
                        httpx.OkJsonCtx(r.Context(), w, resp)
                }
        } {{else}}return func(w http.ResponseWriter, r *http.Request) {
                l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx, r, w)
                err := l.{{.Call}}()
                if err != nil {
                        httpx.ErrorCtx(r.Context(), w, err)
                } else {
                        httpx.Ok(w)
                }
        }{{end}}
}