package errors

import (
	"context"

	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/utils"
)

// UnaryServerInterceptor returns a new unary server interceptor that unwraps errors before returning them
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		res, err := handler(ctx, req)
		if err != nil {
			return res, utils.UnwrapError(ctx, err)
		}
		return res, err
	}
}

// StreamServerInterceptor returns a new stream server interceptors that unwraps errors before returning them
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return utils.UnwrapError(stream.Context(), handler(srv, stream))
	}
}
