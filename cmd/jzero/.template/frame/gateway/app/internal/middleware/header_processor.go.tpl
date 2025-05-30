package middleware

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/gateway"
	"google.golang.org/grpc"
)

func WithHeaderProcessor() gateway.Option {
	return gateway.WithHeaderProcessor(func(header http.Header) []string {
		var headers []string
		//// You can add header from request header here
		//// for example
		//for k, v := range header {
		//	if k == "Authorization" {
		//		headers = append(headers, fmt.Sprintf("%s:%s", k, strings.Join(v, ";")))
		//	}
		//}
		return headers
	})
}

func WithValueMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	//md, b := metadata.FromIncomingContext(ctx)
	//if !b {
	//	return handler(ctx, req)
	//}
	//// You can verify Authorization here and set user info in context value
	//// get Authorization
	//value := md.Get("Authorization")
	//if len(value) == 1 {
	//	// set context value
	//	ctx = context.WithValue(ctx, "Authorization", value[0])
	//}
	return handler(ctx, req)
}
