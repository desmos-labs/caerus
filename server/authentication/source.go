package authentication

import (
	"github.com/desmos-labs/caerus/types"
)

// Source represents the interface that must be implemented in order to get the encrypted session of a user
type Source interface {
	// GetSession returns the session associated to the given token
	// or an error if something goes wrong
	GetSession(token string) (*types.EncryptedSession, error)
}
