package authentication

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unAuthenticatedContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, DataKey, &UnAuthenticatedContext{})
}

func authenticatedUserContext(ctx context.Context, data *AuthenticatedUserData) context.Context {
	return context.WithValue(ctx, DataKey, data)
}

func authenticatedAppContext(ctx context.Context, data *AuthenticatedAppData) context.Context {
	return context.WithValue(ctx, DataKey, data)
}

// BearerTokenAuthFunction returns a function that checks whether the given context.Context
// object contains the metadata related to authentication.
// Based on the metadata supplied within the context, it returns either one of the three context
// implementations: UnAuthenticatedContext, AuthenticatedUserData or AuthenticatedAppData.
func BearerTokenAuthFunction(source Source) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			if errStat, ok := status.FromError(err); ok && errStat.Code() == codes.Unauthenticated {
				return unAuthenticatedContext(ctx), nil
			}
			return nil, err
		}

		// Try authenticating as a user
		session, err := source.GetUserSession(token)
		if err != nil {
			return nil, err
		}

		if session != nil {
			return authenticatedUserContext(ctx, &AuthenticatedUserData{
				Token:         token,
				DesmosAddress: session.DesmosAddress,
			}), nil
		}

		// Try authenticating as an application
		appToken, err := source.GetAppToken(token)
		if err != nil {
			return nil, err
		}

		if appToken != nil {
			return authenticatedAppContext(ctx, &AuthenticatedAppData{
				Token: token,
				AppID: appToken.AppID,
			}), nil
		}

		// Return unauthenticated context
		return unAuthenticatedContext(ctx), nil
	}
}
