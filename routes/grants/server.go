package grants

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/posthog/posthog-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/desmos-labs/caerus/analytics"
	"github.com/desmos-labs/caerus/authentication"
)

var (
	_ GrantsServiceServer = &Server{}
)

type Server struct {
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func NewServerFromEnvVariables(chainClient ChainClient, cdc codec.Codec, db Database) *Server {
	return NewServer(
		NewHandler(chainClient, cdc, db),
	)
}

func (s *Server) RequestFeeAllowance(ctx context.Context, request *RequestFeeAllowanceRequest) (*emptypb.Empty, error) {
	appData, err := authentication.GetAuthenticatedAppData(ctx)
	if err != nil {
		return nil, err
	}

	// Handle the request
	req := NewRequestFeeGrantRequest(appData.AppID, request.UserDesmosAddress, request.Allowance)
	err = s.handler.HandleFeeGrantRequest(req)
	if err != nil {
		return nil, err
	}

	// Log the event
	analytics.Enqueue(posthog.Capture{
		DistinctId: appData.AppID,
		Event:      "Requested fee grant",
		Properties: posthog.NewProperties().
			Set(analytics.KeyUserAddress, request.UserDesmosAddress),
	})

	return &emptypb.Empty{}, nil
}
