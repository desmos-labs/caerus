package notifications

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/types"
)

var (
	_ NotificationsServiceServer = &Server{}
)

type Server struct {
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func NewServerFromEnvVariables(firebase Firebase, db Database) *Server {
	return NewServer(
		NewHandler(firebase, db),
	)
}

func (s Server) SendNotification(ctx context.Context, request *SendNotificationRequest) (*emptypb.Empty, error) {
	appData, err := authentication.GetAuthenticatedAppData(ctx)
	if err != nil {
		return nil, err
	}

	// Parse the notification data
	var notification types.Notification
	err = json.Unmarshal(request.Notification, &notification)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid notification data")
	}

	// Build and validate the request
	req := NewSendAppNotificationRequest(appData.AppID, request.UserAddresses, &notification)
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Handle the request
	err = s.handler.HandleSendNotificationRequest(req)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, err
}
