package grants

import (
	"context"

	"github.com/posthog/posthog-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/desmos-labs/caerus/analytics"
	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/utils"
)

var (
	_ GrantsServiceServer = &Server{}
)

type Server struct {
	handler *Handler
}

func NewServer(client ChainClient, db Database) *Server {
	return &Server{
		handler: NewHandler(client, db),
	}
}

func (s *Server) RequestFeeAllowance(ctx context.Context, request *RequestFeeAllowanceRequest) (*emptypb.Empty, error) {
	appCtx, err := authentication.GetAppContext(ctx)
	if err != nil {
		return nil, err
	}

	// Handle the request
	err = s.handler.HandleFeeGrantRequest(NewRequestFeeGrantRequest(appCtx.AppID, request.UserDesmosAddress))
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	// Log the event
	analytics.Enqueue(posthog.Capture{
		DistinctId: appCtx.AppID,
		Event:      "Requested fee grant",
		Properties: posthog.NewProperties().
			Set(analytics.KeyUserAddress, request.UserDesmosAddress),
	})

	return &emptypb.Empty{}, err
}
