package authentication

type ContextAuthenticationKey string

var (
	DataKey ContextAuthenticationKey = "AuthenticationData"
)

// UnAuthenticatedContext represents the gRPC context that will be used when the user is not authenticated
type UnAuthenticatedContext struct {
}

// AuthenticatedUserData contains the data of an authenticated user
type AuthenticatedUserData struct {
	Token         string
	DesmosAddress string
}

// AuthenticatedAppData contains the data of an authenticated application
type AuthenticatedAppData struct {
	Token string
	AppID string
}
