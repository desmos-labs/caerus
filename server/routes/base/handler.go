package base

import (
	"net/http"

	"github.com/desmos-labs/caerus/server/utils"
	"github.com/desmos-labs/caerus/types"
)

// Handler represents a basic handler that provides basic functionalities
type Handler struct {
	db Database
}

func NewHandler(db Database) *Handler {
	return &Handler{
		db: db,
	}
}

// GetSession returns
func (h *Handler) GetSession(token string) (*types.EncryptedSession, error) {
	// Check the session validity
	session, err := h.db.GetSession(token)
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
