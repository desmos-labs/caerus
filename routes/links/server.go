package links

import (
	"context"

	"github.com/desmos-labs/caerus/authentication"
)

var (
	_ LinksServiceServer = &Server{}
)

type Server struct {
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	return &Server{
		handler: handler,
	}
}

// NewServerFromEnvVariables builds a new Server instance reading the config from the env variables
func NewServerFromEnvVariables(client DeepLinksClient, db Database) *Server {
	return NewServer(NewHandlerFromEnvVariables(client, db))
}

// --------------------------------------------------------------------------------------------------------------------

// CreateLink implements LinksServiceServer
func (s *Server) CreateLink(ctx context.Context, request *CreateLinkRequest) (*CreateLinkResponse, error) {
	// Get the app information
	app, err := authentication.GetAuthenticatedAppData(ctx)
	if err != nil {
		return nil, err
	}

	// Build and handle the request
	req := NewGenerateDeepLinkRequest(app.AppID, request.LinkConfiguration, request.ApiKey)
	return s.handler.HandleGenerateDeepLinkRequest(req)
}
