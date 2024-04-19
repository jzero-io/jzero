package middlewares

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var RateLimit syncx.Limit

func GrpcRateLimit(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if RateLimit.TryBorrow() {
		defer func() {
			if err := RateLimit.Return(); err != nil {
				logx.Error(err)
			}
		}()
		return handler(ctx, req)
	} else {
		logx.Errorf("concurrent connections over %d, rejected with code %d",
			100, http.StatusServiceUnavailable)
		return nil, status.Error(codes.Unavailable, "concurrent connections over limit")
	}
}
