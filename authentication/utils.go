package authentication

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetAuthenticatedUserData(ctx context.Context) (*AuthenticatedUserData, error) {
	userContext, ok := ctx.Value(DataKey).(*AuthenticatedUserData)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	return userContext, nil
}

func GetAuthenticatedAppData(ctx context.Context) (*AuthenticatedAppData, error) {
	appContext, ok := ctx.Value(DataKey).(*AuthenticatedAppData)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	return appContext, nil
}
