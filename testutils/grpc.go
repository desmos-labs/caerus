package testutils

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
	netListener, err := net.Listen("tcp", ":19090")
	if err != nil {
		return nil, err
	}
	go server.Serve(netListener)

	// Create the connection
	return grpc.Dial("localhost:19090", grpc.WithTransportCredentials(insecure.NewCredentials()))
}
