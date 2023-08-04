package grants

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	SaveFeeGrantRequest(request types.FeeGrantRequest) error
	HasFeeGrantBeenGrantedToUser(user string) (bool, error)
}
