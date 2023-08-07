package grants

import (
	"github.com/desmos-labs/caerus/server/routes/base"
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	base.Database

	GetApp(appID string) (*types.Application, bool, error)

	GetAppFeeGrantRequestsLimit(appID string) (uint64, error)
	GetAppFeeGrantRequestsCount(appID string) (uint64, error)

	SaveFeeGrantRequest(request types.FeeGrantRequest) error
	HasFeeGrantBeenGrantedToUser(appID string, user string) (bool, error)
}
