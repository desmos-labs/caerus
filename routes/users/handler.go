package users

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/desmos-labs/caerus/types"
	serverutils "github.com/desmos-labs/caerus/utils"
)

type Handler struct {
	cdc   codec.Codec
	amino *codec.LegacyAmino
	db    Database
}

// NewHandler returns a new Handler instance
func NewHandler(cdc codec.Codec, amino *codec.LegacyAmino, db Database) *Handler {
	return &Handler{
		db:    db,
		cdc:   cdc,
		amino: amino,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// HandleNonceRequest returns the proper nonce for the given request
func (h *Handler) HandleNonceRequest(request *GetNonceRequest) (*GetNonceResponse, error) {
	_, err := sdk.AccAddressFromBech32(request.UserDesmosAddress)
	if err != nil {
		return nil, serverutils.WrapErr(http.StatusBadRequest, "invalid Desmos address")
	}

	// Get the nonce
	nonce, err := types.CreateNonce(request.UserDesmosAddress)
	if err != nil {
		return nil, err
	}

	// Store the nonce
	err = h.db.SaveNonce(nonce)
	if err != nil {
		return nil, err
	}

	return NewGetNonceResponse(nonce.Value), nil
}

// HandleAuthenticationRequest checks the given request to make sure the user has authenticated correctly
func (h *Handler) HandleAuthenticationRequest(request *types.SignedRequest) (*LoginResponse, error) {
	// Verify the request
	memo, err := request.Verify(h.cdc, h.amino)
	if err != nil {
		return nil, err
	}

	// Get the nonce from the database
	nonce, err := h.db.GetNonce(request.DesmosAddress, memo)
	if err != nil {
		return nil, err
	}

	if nonce == nil {
		return nil, serverutils.WrapErr(http.StatusBadRequest, "wrong nonce value")
	}

	// Remove the nonce
	err = h.db.DeleteNonce(nonce)
	if err != nil {
		return nil, err
	}

	// Verify the validity of the nonce
	if nonce.IsExpired() {
		return nil, serverutils.WrapErr(http.StatusBadRequest, "nonce is expired")
	}

	// Create the session
	session, err := types.CreateUserSession(request.DesmosAddress)
	if err != nil {
		return nil, err
	}

	// Save the session inside the database
	err = h.db.SaveSession(session)
	if err != nil {
		return nil, err
	}

	// Save the login information
	err = h.db.UpdateLoginInfo(request.DesmosAddress)
	if err != nil {
		return nil, err
	}

	return NewLoginResponse(session.Token, &session.ExpiryTime), nil
}

// --------------------------------------------------------------------------------------------------------------------

// HandleRefreshSessionRequest refreshes the session associated with the given request
func (h *Handler) HandleRefreshSessionRequest(token string) (*RefreshSessionResponse, error) {
	// Check the session validity
	session, err := h.db.GetUserSession(token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, serverutils.WrapErr(http.StatusUnauthorized, "invalid token")
	}

	_, shouldDelete, err := session.Validate()
	if shouldDelete {
		err := h.db.DeleteSession(session.EncryptedToken)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, serverutils.WrapErr(http.StatusUnauthorized, err.Error())
	}

	// Refresh the session
	session = session.Refresh()

	// Update the session
	err = h.db.UpdateSession(session)
	if err != nil {
		return nil, err
	}

	// Save the user information
	err = h.db.SaveUser(session.DesmosAddress)
	if err != nil {
		return nil, err
	}

	return NewSessionRefreshResponse(token, session.ExpiryTime), nil
}

// --------------------------------------------------------------------------------------------------------------------

// HandleLogoutRequest allows to handle the request of logging out a user that is connected using the given session
func (h *Handler) HandleLogoutRequest(req *LogoutUserRequest) error {
	session, err := h.db.GetUserSession(req.Token)
	if err != nil {
		return err
	}

	if req.LogoutAllDevices {
		return h.db.DeleteAllSessions(session.DesmosAddress)
	} else {
		return h.db.DeleteSession(session.EncryptedToken)
	}
}

// HandleDeleteAccountRequest allows to handle the request of deleting a user account
func (h *Handler) HandleDeleteAccountRequest(req *DeleteAccountRequest) error {
	return h.db.DeleteAllSessions(req.UserAddress)
}

// --------------------------------------------------------------------------------------------------------------------

// getDefaultHasuraResponse returns the default Hasura session response
func (h *Handler) getDefaultHasuraResponse() *UnauthorizedSessionResponse {
	return NewUnauthorizedHasuraSessionResponse()
}

// HandleHasuraSessionRequest returns the session used to authenticate a Hasura user
func (h *Handler) HandleHasuraSessionRequest(token string) (HasuraSessionResponse, error) {
	session, err := h.db.GetUserSession(token)
	if err != nil {
		return nil, err
	}

	defaultResponse := h.getDefaultHasuraResponse()

	return NewAuthorizedHasuraSessionResponse(defaultResponse, session.DesmosAddress), nil
}

// GetUnauthorizedHasuraSession returns the session used to authenticate an unauthorized Hasura user
func (h *Handler) GetUnauthorizedHasuraSession() HasuraSessionResponse {
	return h.getDefaultHasuraResponse()
}

// --------------------------------------------------------------------------------------------------------------------

// HandleRegisterUserDeviceTokenRequest handles the request to register a new device token
func (h *Handler) HandleRegisterUserDeviceTokenRequest(req *RegisterUserDeviceTokenRequest) error {
	return h.db.SaveUserNotificationDeviceToken(types.NewUserNotificationDeviceToken(req.UserAddress, req.DeviceToken))
}
