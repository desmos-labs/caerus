package server

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/logging"
	"github.com/desmos-labs/caerus/server/errors"

	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// DefaultInterceptors returns the default interceptors that should be used by the server
func DefaultInterceptors(source authentication.Source) []grpc.ServerOption {
	loggingFunction := logging.ZeroLogInterceptorLogger(log.Logger)
	authenticationFunction := authentication.BearerTokenAuthFunction(source)

	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpclogging.UnaryServerInterceptor(loggingFunction),
			grpcauth.UnaryServerInterceptor(authenticationFunction),
			errors.UnaryServerInterceptor(),
		),

		grpc.ChainStreamInterceptor(
			grpclogging.StreamServerInterceptor(loggingFunction),
			grpcauth.StreamServerInterceptor(authenticationFunction),
			errors.StreamServerInterceptor(),
		),
	}
}
