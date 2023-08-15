package notifications

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

var (
	_ NotificationsServiceServer = &Server{}
)

type Server struct {
	handler *Handler
}

func NewServer(firebase Firebase, db Database) *Server {
	return &Server{
		handler: NewHandler(firebase, db),
	}
}

func (s Server) SendNotification(ctx context.Context, request *SendNotificationRequest) (*emptypb.Empty, error) {
	appCtx, err := authentication.GetAppContext(ctx)
	if err != nil {
		return nil, err
	}

	// Parse the notification data
	var notification types.Notification
	err = json.Unmarshal(request.Notification, &notification)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid notification")
	}

	// Handle the request
	req := NewSendAppNotificationRequest(appCtx.AppID, request.DeviceTokens, &notification)
	err = s.handler.HandleSendNotificationRequest(req)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	return &emptypb.Empty{}, err
}
