package operations

import (
	"fmt"
	"strings"

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
		err := GrantAuthorizations(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error while granting authorizations")
		}
	})
}

// GrantAuthorizations grants the authorizations to the users that have requested them
func GrantAuthorizations(ctx scheduler.Context) error {
	// Get at most 100 authorizations requests
	feeGrantRequests, err := ctx.Database.GetNotGrantedFeeGrantRequests(100)
	if err != nil {
		return err
	}

	// If there are no fee grant requests, return
	if len(feeGrantRequests) == 0 {
		return nil
	}

	// Group the authorizations based on the application that is associated to them
	appsFeeGrantRequests := map[string][]types.FeeGrantRequest{}
	for _, feeGrantRequest := range feeGrantRequests {
		appsFeeGrantRequests[feeGrantRequest.AppID] = append(appsFeeGrantRequests[feeGrantRequest.AppID], feeGrantRequest)
	}

	// Get the grant requests ids.
	// This will be used later on to mark the fee requests as granted
	var grantedGrantRequestsIDs []string

	// Get the list of users to whom the fee allowances have been granted.
	// This will be used later on to notify applications about which users the grants have been given to
	grantedUsers := map[string][]string{}

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

			// Skip this application since we are not going to be able to properly execute the message
			continue
		}

		// Check if the app still has the on-chain MsgExec authorization
		// If it does not have an authorization, send a notification
		hasOnChainAuthorization, err := ctx.ChainClient.HasGrantedMsgGrantAllowanceAuthorization(app.WalletAddress)
		if err != nil {
			return err
		}

		if !hasOnChainAuthorization {
			err = ctx.FirebaseClient.SendNotificationToApp(appID, types.Notification{
				Data: map[string]string{
					types.NotificationTypeKey:    "missing_authorization",
					types.NotificationMessageKey: "Your application is missing the on-chain authorization to be able to send fee allowances",
				},
			})
			if err != nil {
				return err
			}

			// Skip this application since we are not going to be able to properly execute the message
			continue
		}

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

			grantedGrantRequestsIDs = append(grantedGrantRequestsIDs, grantRequest.ID)
			grantedUsers[appID] = append(grantedUsers[appID], grantRequest.DesmosAddress)

			// Only add the message if the fee allowance has not been granted yet
			hasFeeGrant, err := ctx.ChainClient.HasFeeGrant(feeGrantGranteeAddress.String(), feeGrantGranterAddress.String())
			if err != nil {
				return err
			}

			if !hasFeeGrant {
				grantAllowanceMsgs = append(grantAllowanceMsgs, msgGrantAllowance)
			}
		}

		if len(grantAllowanceMsgs) == 0 {
			continue
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

	// If there are no messages, just return
	if len(msgExecMsgs) == 0 {
		return nil
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
	err = ctx.Database.SetFeeGrantRequestsGranted(grantedGrantRequestsIDs)
	if err != nil {
		return err
	}

	// Send a notification to the applications
	for appID, usersAddresses := range grantedUsers {
		err = ctx.FirebaseClient.SendNotificationToApp(appID, types.Notification{
			Data: map[string]string{
				types.NotificationTypeKey:    "fee_allowances_granted",
				types.NotificationMessageKey: "Your fee allowances have been granted",
				"users":                      strings.Join(usersAddresses, ","),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
