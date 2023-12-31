package users

import (
	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	authentication.Database

	GetNonce(desmosAddress string, value string) (*types.EncryptedNonce, error)
	SaveNonce(nonce *types.Nonce) error
	DeleteNonce(nonce *types.EncryptedNonce) error

	SaveUser(desmosAddress string) error
	UpdateLoginInfo(desmosAddress string) error
	DeleteAllSessions(desmosAddress string) error

	SaveUserNotificationDeviceToken(token *types.UserNotificationDeviceToken) error
}
