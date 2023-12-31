package users

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	AuthenticatedUserRole = "user"
	UnauthorizedUserRole  = "anonymous"
)

// NewGetNonceResponse builds a new GetNonceRequest instance
func NewGetNonceResponse(nonce string) *GetNonceResponse {
	return &GetNonceResponse{
		Nonce: nonce,
	}
}

// NewLoginResponse builds a new LoginResponse instance
func NewLoginResponse(token string, expirationTime *time.Time) *LoginResponse {
	return &LoginResponse{
		Token:          token,
		ExpirationTime: expirationTime,
	}
}

// NewSessionRefreshResponse builds a new SessionRefreshResponse instance
func NewSessionRefreshResponse(token string, expirationTime *time.Time) *RefreshSessionResponse {
	return &RefreshSessionResponse{
		Token:          token,
		ExpirationTime: expirationTime,
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
	DeviceToken string
}

func NewRegisterUserDeviceTokenRequest(userAddress string, token string) *RegisterUserDeviceTokenRequest {
	return &RegisterUserDeviceTokenRequest{
		UserAddress: userAddress,
		DeviceToken: token,
	}
}

func (r RegisterUserDeviceTokenRequest) Validate() error {
	if r.DeviceToken == "" {
		return fmt.Errorf("invalid device token")
	}

	return nil
}
