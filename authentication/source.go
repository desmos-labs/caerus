package authentication

import (
	"net/http"

	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

// Source represents the interface that must be implemented in order to get the encrypted session of a user
type Source interface {
	// GetUserSession returns the user session associated to the given token
	// or an error if something goes wrong
	GetUserSession(token string) (*types.EncryptedUserSession, error)

	// GetAppToken returns the app session associated to the given token
	// or an error if something goes wrong
	GetAppToken(token string) (*types.EncryptedAppToken, error)
}

// Database represents the interface that must be implemented in order to properly read the data required during
// the authentication processes
type Database interface {
	SaveSession(session *types.UserSession) error
	GetUserSession(token string) (*types.EncryptedUserSession, error)
	UpdateSession(session *types.EncryptedUserSession) error
	DeleteSession(token string) error

	GetAppToken(token string) (*types.EncryptedAppToken, error)
}

// -------------------------------------------------------------------------------------------------------------------

var (
	_ Source = &BaseAuthSource{}
)

// BaseAuthSource represents a basic implementation for the Source interface
type BaseAuthSource struct {
	db Database
}

func NewBaseAuthSource(db Database) *BaseAuthSource {
	return &BaseAuthSource{
		db: db,
	}
}

// GetUserSession implements authentication.Source
func (h *BaseAuthSource) GetUserSession(token string) (*types.EncryptedUserSession, error) {
	// Check the session validity
	session, err := h.db.GetUserSession(token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, utils.WrapErr(http.StatusUnauthorized, "invalid token")
	}

	shouldRefresh, shouldDelete, err := session.Validate()
	if shouldDelete {
		err := h.db.DeleteSession(session.EncryptedToken)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, utils.WrapErr(http.StatusUnauthorized, err.Error())
	}

	if shouldRefresh {
		session = session.Refresh()
		err = h.db.UpdateSession(session)
		if err != nil {
			return nil, err
		}
	}

	return session, nil
}

// GetAppToken implements authentication.Source
func (h *BaseAuthSource) GetAppToken(token string) (*types.EncryptedAppToken, error) {
	// Check the session validity
	encryptedToken, err := h.db.GetAppToken(token)
	if err != nil {
		return nil, err
	}

	if encryptedToken == nil {
		return nil, utils.WrapErr(http.StatusUnauthorized, "invalid token")
	}

	return encryptedToken, nil
}
