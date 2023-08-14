package authentication

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetUserContext(ctx context.Context) (*UserAuthenticatedContext, error) {
	userContext, ok := ctx.(*UserAuthenticatedContext)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	return userContext, nil
}

func GetAppContext(ctx context.Context) (*AppAuthenticatedContext, error) {
	appContext, ok := ctx.(*AppAuthenticatedContext)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	return appContext, nil
}
