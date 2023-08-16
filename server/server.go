package server

import (
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/logging"
)

// JoinServerOptions joins the given options as a single grpc.ServerOption list
func JoinServerOptions(options ...[]grpc.ServerOption) []grpc.ServerOption {
	var result []grpc.ServerOption
	for _, option := range options {
		result = append(result, option...)
	}
	return result
}

// New returns a new server gRPC instance configured as it should be
func New(db authentication.Database) *grpc.Server {
	return grpc.NewServer(JoinServerOptions(
		logging.NewLoggingInterceptor(),
		authentication.NewAuthInterceptors(authentication.NewBaseAuthSource(db)),
	)...)
}
