package grants

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	GetApp(appID string) (*types.Application, bool, error)

	GetAppFeeGrantRequestsLimit(appID string) (uint64, error)
	GetAppFeeGrantRequestsCount(appID string) (uint64, error)

	SaveFeeGrantRequest(request types.FeeGrantRequest) error
	HasFeeGrantBeenGrantedToUser(appID string, user string) (bool, error)

	GetFeeGrantRequest(appID string, requestID string) (*types.FeeGrantRequest, bool, error)
}

type ChainClient interface {
	HasGrantedMsgGrantAllowanceAuthorization(appAddress string) (bool, error)

	HasFunds(address string) (bool, error)
	HasFeeGrant(userAddress string, granterAddress string) (bool, error)
}
