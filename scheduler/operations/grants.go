package operations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/caerus/scheduler"
	"github.com/desmos-labs/caerus/types"
)

func RegisterGrantsOperations(ctx scheduler.Context, scheduler *gocron.Scheduler) {
	scheduler.Every(10).Seconds().Do(func() {
		err := grantAuthorizations(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error while granting authorizations")
		}
	})
}

// grantAuthorizations grants the authorizations to the users that have requested them
func grantAuthorizations(ctx scheduler.Context) error {
	// Get at most 100 authorizations requests
	feeGrantRequests, err := ctx.Database.GetNotGrantedFeeGrantRequests(100)
	if err != nil {
		return err
	}

	// Group the authorizations based on the application that is associated to them
	appsFeeGrantRequests := map[string][]types.FeeGrantRequest{}
	for _, feeGrantRequest := range feeGrantRequests {
		appsFeeGrantRequests[feeGrantRequest.AppID] = append(appsFeeGrantRequests[feeGrantRequest.AppID], feeGrantRequest)
	}

	// Get the grant requests ids
	var grantRequestsIDs []string

	// Build the messages to be sent
	var msgExecMsgs []sdk.Msg
	for appID, appGrantRequests := range appsFeeGrantRequests {
		// Get the application wallet address
		app, found, err := ctx.Database.GetApp(appID)
		if err != nil {
			return err
		}

		if !found {
			log.Error().Str("application id", appID).Msg("application not found while trying to grant a fee allowance")
			continue
		}

		// TODO: Check if the app still has the on-chain MsgExec authorization
		// If it does not have an authorization, send a notification

		// Build the MsgExec instances
		var grantAllowanceMsgs []sdk.Msg
		for _, grantRequest := range appGrantRequests {
			// Parse the addresses
			feeGrantGranterAddress, err := sdk.AccAddressFromBech32(app.WalletAddress)
			if err != nil {
				return err
			}

			feeGrantGranteeAddress, err := sdk.AccAddressFromBech32(grantRequest.DesmosAddress)
			if err != nil {
				return err
			}

			msgGrantAllowance, err := feegrant.NewMsgGrantAllowance(grantRequest.Allowance, feeGrantGranterAddress, feeGrantGranteeAddress)
			if err != nil {
				return err
			}

			grantRequestsIDs = append(grantRequestsIDs, grantRequest.ID)
			grantAllowanceMsgs = append(grantAllowanceMsgs, msgGrantAllowance)
		}

		// Parse the messages
		authzGranteeAddress, err := sdk.AccAddressFromBech32(ctx.ChainClient.AccAddress())
		if err != nil {
			return err
		}

		// Build the MsgExec instance
		msgExec := authz.NewMsgExec(authzGranteeAddress, grantAllowanceMsgs)

		msgExecMsgs = append(msgExecMsgs, &msgExec)
	}

	// Broadcast the transaction
	res, err := ctx.ChainClient.BroadcastTxCommit(&wallettypes.TransactionData{
		Messages: msgExecMsgs,
		Memo:     "Sent from Caerus",
		GasAuto:  true,
		FeeAuto:  true,
	})
	if err != nil {
		return err
	}

	if res.Code != 0 {
		return fmt.Errorf("error while broadcasting fee allowance transactions: %s", res.RawLog)
	}

	// Set the grant requests as granted
	err = ctx.Database.SetFeeGrantRequestsGranted(grantRequestsIDs)
	if err != nil {
		return err
	}

	// TODO: Send a notification to the applications

	return nil
}
