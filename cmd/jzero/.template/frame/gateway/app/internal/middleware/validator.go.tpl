package middleware

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Validator struct {
	v protovalidate.Validator
}

func NewValidator() *Validator {
	v, err := protovalidate.New()
	logx.Must(err)
	return &Validator{v: v}
}

func (v *Validator) UnaryServerMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		switch req.(type) {
		case proto.Message:
			if err := v.v.Validate(req.(proto.Message)); err != nil {
				var valErr *protovalidate.ValidationError
				if ok := errors.As(err, &valErr); ok && len(valErr.ToProto().GetViolations()) > 0 {
					return nil, status.Error(codes.InvalidArgument, valErr.ToProto().GetViolations()[0].GetMessage())
				}
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}
		return handler(ctx, req)
	}
}
