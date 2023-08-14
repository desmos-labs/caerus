package authentication

import (
	"github.com/desmos-labs/caerus/types"
)

// Source represents the interface that must be implemented in order to get the encrypted session of a user
type Source interface {
	// GetUserSession returns the user session associated to the given token
	// or an error if something goes wrong
	GetUserSession(token string) (*types.EncryptedUserSession, error)

	// GetAppSession returns the app session associated to the given token
	// or an error if something goes wrong
	GetAppToken(token string) (*types.EncryptedAppToken, error)
}
