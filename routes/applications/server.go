package applications

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/desmos-labs/caerus/authentication"
)

var (
	_ ApplicationServiceServer = &ApplicationServer{}
)

type ApplicationServer struct {
	handler *Handler
}

func NewServer(handler *Handler) *ApplicationServer {
	return &ApplicationServer{
		handler: handler,
	}
}

func NewServerFromEnvVariables(db Database) *ApplicationServer {
	return NewServer(
		NewHandler(db),
	)
}

func (a ApplicationServer) RegisterNotificationToken(ctx context.Context, request *RegisterNotificationTokenRequest) (*emptypb.Empty, error) {
	appData, err := authentication.GetAuthenticatedAppData(ctx)
	if err != nil {
		return nil, err
	}

	err = a.handler.HandleRegisterAppDeviceTokenRequest(NewRegisterAppDeviceTokenRequest(appData.AppID, request.Token))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a ApplicationServer) DeleteApp(ctx context.Context, request *DeleteAppRequest) (*emptypb.Empty, error) {
	userData, err := authentication.GetAuthenticatedUserData(ctx)
	if err != nil {
		return nil, err
	}

	// Handle the request
	req := NewDeleteApplicationRequest(userData.DesmosAddress, request.AppId)
	err = a.handler.HandleDeleteApplicationRequest(req)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
