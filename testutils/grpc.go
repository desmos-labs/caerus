package testutils

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/server"
)

// CreateServer creates a GRPC server instance with the proper options
func CreateServer(db authentication.Database) *grpc.Server {
	return server.New(db)
}

// StartServerAndConnect starts the given server and connects to it, returning the gRPC connection
// instance that can be used to create clients
func StartServerAndConnect(server *grpc.Server) (*grpc.ClientConn, error) {
	// Start the server
	//nolint:gosec // It is wanted to bind to all interfaces since it's just a test server
	netListener := bufconn.Listen(1028 * 1024)
	go server.Serve(netListener)

	// Create the connection
	return grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return netListener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
}
