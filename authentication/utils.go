package authentication

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnauthenticated = status.Error(codes.Unauthenticated, "Unauthenticated")
)

// AuthenticateUserOrApp checks if the given context contains either an authenticated user or an authenticated app.
// If so, it returns nil, otherwise it returns an error.
func AuthenticateUserOrApp(ctx context.Context) error {
	_, ok := ctx.Value(DataKey).(*AuthenticatedUserData)
	if ok {
		return nil
	}

	_, ok = ctx.Value(DataKey).(*AuthenticatedAppData)
	if ok {
		return nil
	}

	return ErrUnauthenticated
}

// GetAuthenticatedUserData returns the authenticated user data from the given context.
// If the user is not authenticated, an error is returned instead.
func GetAuthenticatedUserData(ctx context.Context) (*AuthenticatedUserData, error) {
	userContext, ok := ctx.Value(DataKey).(*AuthenticatedUserData)
	if !ok {
		return nil, ErrUnauthenticated
	}
	return userContext, nil
}

// GetAuthenticatedAppData returns the authenticated app data from the given context.
// If the app is not authenticated, an error is returned instead.
func GetAuthenticatedAppData(ctx context.Context) (*AuthenticatedAppData, error) {
	appContext, ok := ctx.Value(DataKey).(*AuthenticatedAppData)
	if !ok {
		return nil, ErrUnauthenticated
	}
	return appContext, nil
}
