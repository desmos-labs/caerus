package user

import (
	"github.com/desmos-labs/caerus/server/routes/base"
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	base.Database
	GetNonce(desmosAddress string, value string) (*types.EncryptedNonce, error)
	SaveNonce(nonce *types.Nonce) error
	DeleteNonce(nonce *types.EncryptedNonce) error

	SaveUser(desmosAddress string) error
	UpdateLoginInfo(desmosAddress string) error
	DeleteAllSessions(desmosAddress string) error

	SaveFeeGrantRequest(request types.FeeGrantRequest) error
	HasFeeGrantBeenGrantedToUser(user string) (bool, error)
}
