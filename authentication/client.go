package authentication

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

// SetupContextWithAuthorization setups the given context to properly include the given token as
// the authorization header's value
func SetupContextWithAuthorization(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", token))
}
