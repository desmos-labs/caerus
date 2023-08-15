package users

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/posthog/posthog-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/desmos-labs/caerus/analytics"
	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

var (
	_ UsersServiceServer = &Server{}
)

type Server struct {
	handler *Handler
}

func NewServer(cdc codec.Codec, amino *codec.LegacyAmino, db Database) *Server {
	return &Server{
		handler: NewHandler(cdc, amino, db),
	}
}

func (s *Server) GetNonce(ctx context.Context, request *GetNonceRequest) (*GetNonceResponse, error) {
	res, err := s.handler.HandleNonceRequest(request)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	// Log the event
	analytics.Enqueue(posthog.Capture{
		DistinctId: request.UserDesmosAddress,
		Event:      "Requested Nonce",
	})

	return res, nil
}

func (s *Server) Login(ctx context.Context, request *types.SignedRequest) (*LoginResponse, error) {
	res, err := s.handler.HandleAuthenticationRequest(request)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	// Log the request
	analytics.Enqueue(posthog.Capture{
		DistinctId: request.DesmosAddress,
		Event:      "Logged In",
	})

	return res, err
}

func (s *Server) RefreshSession(ctx context.Context, _ *emptypb.Empty) (*RefreshSessionResponse, error) {
	userCtx, err := authentication.GetUserContext(ctx)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	res, err := s.handler.HandleRefreshSessionRequest(userCtx.Token)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	return res, err
}

func (s *Server) RegisterDeviceNotificationToken(ctx context.Context, request *RegisterNotificationDeviceTokenRequest) (*emptypb.Empty, error) {
	userCtx, err := authentication.GetUserContext(ctx)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	// Handle the request
	req := NewRegisterUserDeviceTokenRequest(userCtx.DesmosAddress, request.DeviceToken)
	err = s.handler.HandleRegisterUserDeviceTokenRequest(req)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Logout(ctx context.Context, request *LogoutRequest) (*emptypb.Empty, error) {
	userCtx, err := authentication.GetUserContext(ctx)
	if err != nil {
		return nil, err
	}

	// Handle the request
	err = s.handler.HandleLogoutRequest(NewLogoutUserRequest(userCtx.Token, request.LogoutFromAll))
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	return &emptypb.Empty{}, err
}

func (s *Server) DeleteAccount(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	userCtx, err := authentication.GetUserContext(ctx)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	// Handle the request
	req := NewDeleteAccountRequest(userCtx.DesmosAddress)
	err = s.handler.HandleDeleteAccountRequest(req)
	if err != nil {
		return nil, utils.UnwrapError(ctx, err)
	}

	// Log the event
	analytics.Enqueue(posthog.Capture{
		DistinctId: req.UserAddress,
		Event:      "Started Account Deletion",
	})

	return &emptypb.Empty{}, nil
}
