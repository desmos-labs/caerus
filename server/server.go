package server

import (
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/authentication"
)

// New returns a new server gRPC instance configured as it should be
func New(db authentication.Database) *grpc.Server {
	return grpc.NewServer(
		// Register the default interceptors
		DefaultInterceptors(authentication.NewBaseAuthSource(db))...,
	)
}
