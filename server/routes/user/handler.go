package user

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	client "github.com/desmos-labs/caerus/server/chain"
	"github.com/desmos-labs/caerus/server/routes/base"
	serverutils "github.com/desmos-labs/caerus/server/utils"
	"github.com/desmos-labs/caerus/types"
)

type Handler struct {
	*base.Handler
	cdc    codec.Codec
	amino  *codec.LegacyAmino
	client *client.Client
	db     Database
}

// NewHandler returns a new Handler instance
func NewHandler(db Database, cdc codec.Codec, amino *codec.LegacyAmino, client *client.Client) *Handler {
	return &Handler{
		Handler: base.NewHandler(db),
		db:      db,
		cdc:     cdc,
		amino:   amino,
		client:  client,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// HandleNonceRequest returns the proper nonce for the given request
func (h *Handler) HandleNonceRequest(request *NonceRequest) (*NonceResponse, error) {
	_, err := sdk.AccAddressFromBech32(request.DesmosAddress)
	if err != nil {
		return nil, serverutils.WrapErr(http.StatusBadRequest, "invalid Desmos address")
	}

	// Get the nonce
	nonce, err := types.CreateNonce(request.DesmosAddress)
	if err != nil {
		return nil, err
	}

	// Store the nonce
	err = h.db.SaveNonce(nonce)
	if err != nil {
		return nil, err
	}

	return NewNonceResponse(nonce.Value), nil
}

// ParseAuthenticateRequest parses the given request body into an AuthenticationRequest instace
func (h *Handler) ParseAuthenticateRequest(body []byte) (*AuthenticationRequest, error) {
	var req AuthenticationRequest
	return &req, json.Unmarshal(body, &req)
}

// HandleAuthenticationRequest checks the given request to make sure the user has authenticated correctly
func (h *Handler) HandleAuthenticationRequest(request *AuthenticationRequest) (*AuthenticationResponse, error) {
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
	session, err := types.CreateSession(request.DesmosAddress)
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

	return NewAuthenticationResponse(session.Token), nil
}

// --------------------------------------------------------------------------------------------------------------------

// HandleRefreshSessionRequest refreshes the session associated with the given request
func (h *Handler) HandleRefreshSessionRequest(token string) error {
	// Check the session validity
	session, err := h.db.GetSession(token)
	if err != nil {
		return err
	}

	if session == nil {
		return serverutils.WrapErr(http.StatusUnauthorized, "invalid token")
	}

	_, shouldDelete, err := session.Validate()
	if shouldDelete {
		err := h.db.DeleteSession(session.EncryptedToken)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return serverutils.WrapErr(http.StatusUnauthorized, err.Error())
	}

	// Refresh the session
	err = h.db.UpdateSession(session.Refresh())
	if err != nil {
		return err
	}

	// Save the user information
	err = h.db.SaveUser(session.DesmosAddress)
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// HandleLogoutRequest allows to handle the request of logging out a user that is connected using the given session
func (h *Handler) HandleLogoutRequest(req *LogoutRequest) error {
	session, err := h.db.GetSession(req.Token)
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
	session, err := h.GetSession(token)
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
