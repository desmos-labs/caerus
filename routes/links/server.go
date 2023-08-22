package links

import (
	"context"

	"google.golang.org/grpc/codes"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/utils"
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

	// Build and validate the request
	req := NewGenerateGenericDeepLinkRequest(app.AppID, request.LinkConfiguration, request.ApiKey)
	err = req.Validate()
	if err != nil {
		return nil, utils.WrapErr(codes.InvalidArgument, err.Error())
	}

	return s.handler.HandleGenerateGenericDeepLinkRequest(req)
}

// CreateAddressLink implements LinksServiceServer
func (s *Server) CreateAddressLink(ctx context.Context, request *CreateAddressLinkRequest) (*CreateLinkResponse, error) {
	// Get the app information
	app, err := authentication.GetAuthenticatedAppData(ctx)
	if err != nil {
		return nil, err
	}

	// Build and validate the request
	req := NewGenerateAddressLinkRequest(app.AppID, request.Address, request.Chain)
	err = req.Validate()
	if err != nil {
		return nil, utils.WrapErr(codes.InvalidArgument, err.Error())
	}

	return s.handler.HandleGenerateDeepLinkRequest(req)
}

// CreateViewProfileLink implements LinksServiceServer
func (s *Server) CreateViewProfileLink(ctx context.Context, request *CreateViewProfileLinkRequest) (*CreateLinkResponse, error) {
	// Get the app information
	app, err := authentication.GetAuthenticatedAppData(ctx)
	if err != nil {
		return nil, err
	}

	// Build and validate the request
	req := NewGenerateViewProfileLinkRequest(app.AppID, request.Address, request.Chain)
	err = req.Validate()
	if err != nil {
		return nil, utils.WrapErr(codes.InvalidArgument, err.Error())
	}

	return s.handler.HandleGenerateDeepLinkRequest(req)
}

// CreateSendLink implements LinksServiceServer
func (s *Server) CreateSendLink(ctx context.Context, request *CreateSendLinkRequest) (*CreateLinkResponse, error) {
	// Get the app information
	app, err := authentication.GetAuthenticatedAppData(ctx)
	if err != nil {
		return nil, err
	}

	// Build and validate the request
	req := NewGenerateSendTokensLinkRequest(app.AppID, request.Address, request.Chain, request.Amount)
	err = req.Validate()
	if err != nil {
		return nil, utils.WrapErr(codes.InvalidArgument, err.Error())
	}

	return s.handler.HandleGenerateDeepLinkRequest(req)
}
