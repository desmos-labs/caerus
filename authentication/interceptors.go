package authentication

import (
	"context"

	auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

// BearerTokenAuthFunction returns a function that checks whether the given context.Context
// object contains the metadata related to authentication.
// Based on the metadata supplied within the context, it returns either one of the three context
// implementations: UnAuthenticatedContext, UserAuthenticatedContext or AppAuthenticatedContext.
func BearerTokenAuthFunction(source Source) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		// Try authenticating as a user
		session, err := source.GetUserSession(token)
		if err != nil {
			return nil, err
		}

		if session != nil {
			return &UserAuthenticatedContext{
				Context:       ctx,
				Token:         token,
				DesmosAddress: session.DesmosAddress,
			}, nil
		}

		// Try authenticating as an application
		appToken, err := source.GetAppToken(token)
		if err != nil {
			return nil, err
		}

		if appToken != nil {
			return &AppAuthenticatedContext{
				Context: ctx,
				Token:   appToken.TokenValue,
				AppID:   appToken.AppID,
			}, nil
		}

		// Return unauthenticated context
		return &UnAuthenticatedContext{
			Context: ctx,
		}, nil
	}
}

// NewAuthInterceptors returns a list of grpc.ServerOption that can be used to register interceptors
// inside a gRPC server to make it support client-side authentication properly
func NewAuthInterceptors(source Source) []grpc.ServerOption {
	authFunc := BearerTokenAuthFunction(source)
	return []grpc.ServerOption{
		grpc.StreamInterceptor(auth.StreamServerInterceptor(authFunc)),
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(authFunc)),
	}
}
