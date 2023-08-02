package base

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	SaveSession(session *types.Session) error
	GetSession(token string) (*types.EncryptedSession, error)
	UpdateSession(session *types.EncryptedSession) error
	DeleteSession(token string) error
}
