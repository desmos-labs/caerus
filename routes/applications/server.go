package applications

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/utils"
)

var (
	_ ApplicationServiceServer = &ApplicationServer{}
)

type ApplicationServer struct {
	handler *Handler
}

func NewServer(db Database) *ApplicationServer {
	return &ApplicationServer{
		handler: NewHandler(db),
	}
}

func (a ApplicationServer) RegisterNotificationToken(ctx context.Context, request *RegisterNotificationTokenRequest) (*emptypb.Empty, error) {
	appCtx, err := authentication.GetAppContext(ctx)
	if err != nil {
		return nil, err
	}

	err = a.handler.HandleRegisterAppDeviceTokenRequest(NewRegisterAppDeviceTokenRequest(appCtx.AppID, request.Token))
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}

func (a ApplicationServer) DeleteApp(ctx context.Context, request *DeleteAppRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
