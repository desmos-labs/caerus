package users

import (
	"encoding/json"
)

const (
	AuthenticatedUserRole = "user"
	UnauthorizedUserRole  = "anonymous"
)

func NewGetNonceResponse(nonce string) *GetNonceResponse {
	return &GetNonceResponse{
		Nonce: nonce,
	}
}

// --------------------------------------------------------------------------------------------------------------------

func NewLoginResponse(token string) *LoginResponse {
	return &LoginResponse{
		Token: token,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type LogoutUserRequest struct {
	Token            string
	LogoutAllDevices bool `json:"logout_all_devices"`
}

func NewLogoutUserRequest(token string, logoutAll bool) *LogoutUserRequest {
	return &LogoutUserRequest{
		Token:            token,
		LogoutAllDevices: logoutAll,
	}
}

type DeleteAccountRequest struct {
	UserAddress string
}

func NewDeleteAccountRequest(userAddress string) *DeleteAccountRequest {
	return &DeleteAccountRequest{
		UserAddress: userAddress,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type HasuraSessionResponse interface {
	MarshalJSON() ([]byte, error)
}

type HasuraSessionResponseJSON struct {
	UserRole    string `json:"X-Hasura-Role"`
	UserAddress string `json:"X-Hasura-User-Address,omitempty"`
}

type UnauthorizedSessionResponse struct {
}

func NewUnauthorizedHasuraSessionResponse() *UnauthorizedSessionResponse {
	return &UnauthorizedSessionResponse{}
}

func (r *UnauthorizedSessionResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&HasuraSessionResponseJSON{
		UserRole: UnauthorizedUserRole,
	})
}

type AuthorizedSessionResponse struct {
	*UnauthorizedSessionResponse
	userAddress string
}

func NewAuthorizedHasuraSessionResponse(unauthorizedResponse *UnauthorizedSessionResponse, userAddress string) *AuthorizedSessionResponse {
	return &AuthorizedSessionResponse{
		UnauthorizedSessionResponse: unauthorizedResponse,
		userAddress:                 userAddress,
	}
}

func (r *AuthorizedSessionResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&HasuraSessionResponseJSON{
		UserRole:    AuthenticatedUserRole,
		UserAddress: r.userAddress,
	})
}

// --------------------------------------------------------------------------------------------------------------------

// RegisterUserDeviceTokenRequest represents the request sent when a user wants to register
// a new device token to receive notifications
type RegisterUserDeviceTokenRequest struct {
	UserAddress string
	DeviceToken string `json:"token"`
}

func NewRegisterUserDeviceTokenRequest(userAddress string, token string) *RegisterUserDeviceTokenRequest {
	return &RegisterUserDeviceTokenRequest{
		UserAddress: userAddress,
		DeviceToken: token,
	}
}
