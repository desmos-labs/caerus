package testutils

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/desmos-labs/caerus/authentication"
)

// CreateServer creates a GRPC server instance with the proper options
func CreateServer(db authentication.Database) *grpc.Server {
	return grpc.NewServer(authentication.NewAuthInterceptors(authentication.NewBaseAuthSource(db))...)
}

// StartServerAndConnect starts the given server and connects to it, returning the gRPC connection
// instance that can be used to create clients
func StartServerAndConnect(server *grpc.Server) (*grpc.ClientConn, error) {
	// Start the server
	//nolint:gosec // It is wanted to bind to all interfaces since it's just a test server
	netListener, err := net.Listen("tcp", ":19090")
	if err != nil {
		return nil, err
	}
	go server.Serve(netListener)

	// Create the connection
	return grpc.Dial("localhost:19090", grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func SetupContextWithAuthorization(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", token))
}
