package scheduler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"

	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	GetApp(appID string) (*types.Application, bool, error)
	GetAppNotificationTokens(appID string) ([]string, error)

	GetNotGrantedFeeGrantRequests(limit int) ([]types.FeeGrantRequest, error)
	SetFeeGrantRequestsGranted(ids []string) error
}

type Firebase interface {
	SendNotifications(app *types.Application, deviceTokens []string, notification *types.Notification) error
}

type ChainClient interface {
	AccAddress() string
	HasGrantedMsgGrantAllowanceAuthorization(appAddress string) (bool, error)
	BroadcastTxCommit(data *wallettypes.TransactionData) (*sdk.TxResponse, error)
}
