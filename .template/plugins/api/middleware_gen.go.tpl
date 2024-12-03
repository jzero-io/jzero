// Code generated by jzero. DO NOT EDIT.

package middleware

import (
	"fmt"
	"context"
	"net/http"
    "regexp"

	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var (
	_ = fmt.Sprintf("middleware_gen.go")
	_ = context.Background()
    _ = grpc.SupportPackageIsVersion7
)

func RegisterGen(zrpc *zrpc.RpcServer, gw *gateway.Server) {
	gw.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
            {{range $v := .HttpMiddlewares}}
                if matchRoute("{{$v.Name}}", fmt.Sprintf("%s:%s",r.Method,r.URL.Path)) {
                    next = {{$v.Name | FirstUpper | ToCamel}}Middleware(next)
                }
            {{end}}
			next(w, r)
		}
	})

    {{range $v := .ZrpcMiddlewares}}
        zrpc.AddUnaryInterceptors(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
            if matchRoute("{{$v.Name}}", info.FullMethod) {
                return {{$v.Name | FirstUpper | ToCamel}}Middleware(ctx, req, info, handler)
            }
            return handler(ctx, req)
        })
    {{end}}
}

// Define and compile routes
var routesMap map[string][]*regexp.Regexp

// loadRoute compiles and stores a route pattern.
func loadRoutes(middleware string, patterns ...string) {
	if routesMap == nil {
		routesMap = make(map[string][]*regexp.Regexp)
	}

	re := regexp.MustCompile(`\{[^}]+\}`)
	var routes []*regexp.Regexp
	for _, pattern := range patterns {
		pattern = re.ReplaceAllString(pattern, "([^/]+)")
		compiledPattern := "^" + pattern + "$"
		routes = append(routes, regexp.MustCompile(compiledPattern))
	}
	routesMap[middleware] = routes
}

// matchRoute checks if a path matches any compiled route.
func matchRoute(middleware,path string) bool {
	if routes, ok := routesMap[middleware]; ok {
		for _, route := range routes {
			if route.MatchString(path) {
				return true
			}
		}
	}
	return false
}

func init() {
    {{range $v := .HttpMiddlewares}}
        loadRoutes("{{$v.Name}}",
	    {{range $vv := $v.Routes}}"{{$vv}}",
	    {{end}})
    {{end}}

    {{range $v := .ZrpcMiddlewares}}
        loadRoutes("{{$v.Name}}",
	    {{range $vv := $v.Routes}}"{{$vv}}",
	    {{end}})
    {{end}}
}