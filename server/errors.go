package server

import (
	"context"

	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/utils"
)

func ErrorInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	if err != nil {
		return res, utils.UnwrapError(ctx, err)
	}
	return res, err
}

func NewErrorInterceptor() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(ErrorInterceptor),
	}
}
