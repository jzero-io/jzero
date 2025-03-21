package middleware


import (
	"context"

	"google.golang.org/grpc"
	"github.com/zeromicro/go-zero/core/logx"
)

func {{.Name | FirstUpper}}Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
    if MatchRoute("{{.Name}}", info.FullMethod) {
        // do something before middleware
        logx.WithContext(ctx).Info("enter {{.Name}} before middleware")
    }
    hd, err := handler(ctx, req)
    if MatchRoute("{{.Name}}", info.FullMethod) {
        // do something after middleware
        logx.WithContext(ctx).Info("enter {{.Name}} after middleware")
    }
    return hd, err
}

func {{.Name | FirstUpper}}StreamMiddleware(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    if MatchRoute("{{.Name}}", info.FullMethod) {
        // do something before middleware
        logx.WithContext(ss.Context()).Info("enter stream {{.Name}} before middleware")
    }
    hd := handler(srv, ss)
    if MatchRoute("{{.Name}}", info.FullMethod) {
        // do something after middleware
        logx.WithContext(ss.Context()).Info("enter stream {{.Name}} after middleware")
    }
    return hd
}