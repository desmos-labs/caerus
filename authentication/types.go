package authentication

import (
	"context"
)

var (
	_ context.Context = &UnAuthenticatedContext{}
	_ context.Context = &UserAuthenticatedContext{}
	_ context.Context = &AppAuthenticatedContext{}
)

// UnAuthenticatedContext represents the gRPC context that will be used when the user is not authenticated
type UnAuthenticatedContext struct {
	context.Context
}

// UserAuthenticatedContext represents the gRPC context that will be used if the calls are made as an
// authenticated user using the "Authorization: Bearer <token>" authentication method
type UserAuthenticatedContext struct {
	context.Context

	Token         string
	DesmosAddress string
}

// AppAuthenticatedContext represents the gRPC context that will be used if the calls are made as an
// authenticated application using the "Authorization: Bearer <token>" authentication method
type AppAuthenticatedContext struct {
	context.Context

	Token string
	AppID string
}
