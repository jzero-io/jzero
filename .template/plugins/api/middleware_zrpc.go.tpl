package middleware


import (
	"context"

	"google.golang.org/grpc"
)

func {{.Name | FirstUpper}}Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
    return handler(ctx, req)
}

func {{.Name | FirstUpper}}StreamMiddleware(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    return handler(srv, ss)
}