package user

import (
	"encoding/json"

	"github.com/desmos-labs/caerus/types"
)

const (
	AuthenticatedUserRole = "user"
	UnauthorizedUserRole  = "anonymous"
)

type NonceRequest struct {
	// Address of the user requesting to log in
	DesmosAddress string `json:"desmos_address"`
}

func NewNonceRequest(desmosAddress string) *NonceRequest {
	return &NonceRequest{
		DesmosAddress: desmosAddress,
	}
}

type NonceResponse struct {
	// Nonce that should be returned within a signed transaction's memo
	Nonce string `json:"nonce"`
}

func NewNonceResponse(nonce string) *NonceResponse {
	return &NonceResponse{
		Nonce: nonce,
	}
}

type AuthenticationRequest struct {
	*types.SignedRequest
}

type AuthenticationResponse struct {
	// Token that should be sent at each consequent request
	Token string `json:"token"`
}

func NewAuthenticationResponse(token string) *AuthenticationResponse {
	return &AuthenticationResponse{
		Token: token,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type LogoutRequest struct {
	Token            string
	LogoutAllDevices bool `json:"logout_all_devices"`
}

type DeleteAccountRequest struct {
	UserAddress string
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

type SaveContactRequest struct {
	UserAddress string
	Platform    string `json:"platform"`
	Username    string `json:"username"`
}

type DeleteContactRequest struct {
	UserAddress string
	ContactID   string
}
