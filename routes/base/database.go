package base

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	SaveSession(session *types.UserSession) error
	GetUserSession(token string) (*types.EncryptedUserSession, error)
	UpdateSession(session *types.EncryptedUserSession) error
	DeleteSession(token string) error

	GetAppToken(token string) (*types.EncryptedAppToken, error)
}
