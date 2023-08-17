package grants

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"google.golang.org/grpc/codes"

	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

type Handler struct {
	cdc         codec.Codec
	chainClient ChainClient
	db          Database
}

// NewHandler returns a new Handler instance
func NewHandler(client ChainClient, cdc codec.Codec, db Database) *Handler {
	return &Handler{
		cdc:         cdc,
		chainClient: client,
		db:          db,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// HandleFeeGrantRequest handles the request of a fee grant from the user having the given token
func (h *Handler) HandleFeeGrantRequest(req *RequestFeeGrantRequest) error {
	// Unmarshal the allowance
	var allowance feegrant.FeeAllowanceI
	err := h.cdc.UnpackAny(req.Allowance, &allowance)
	if err != nil {
		return err
	}

	// Get the app details
	app, found, err := h.db.GetApp(req.AppID)
	if err != nil {
		return err
	}

	if !found {
		return utils.WrapErr(codes.FailedPrecondition, "application not found")
	}

	// Check if the app has granted a MsgGrantFeeAllowance permission
	hasGrantedAuthorization, err := h.chainClient.HasGrantedMsgGrantAllowanceAuthorization(app.WalletAddress)
	if err != nil {
		return err
	}

	if !hasGrantedAuthorization {
		return utils.WrapErr(codes.FailedPrecondition, "on-chain authorization not found")
	}

	// Check if the application has reached the number of requests
	requestRateLimit, err := h.db.GetAppFeeGrantRequestsLimit(req.AppID)
	if err != nil {
		return err
	}

	requestsCount, err := h.db.GetAppFeeGrantRequestsCount(req.AppID)
	if err != nil {
		return err
	}

	if requestRateLimit > 0 && requestsCount >= requestRateLimit {
		return utils.NewTooManyRequestsError("number of fee grant requests allowed reached")
	}

	// Check if the user has already been granted the fee grant
	hasBeenGranted, err := h.db.HasFeeGrantBeenGrantedToUser(req.AppID, req.DesmosAddress)
	if err != nil {
		return err
	}

	if hasBeenGranted {
		return utils.WrapErr(codes.FailedPrecondition, "you have already been granted the authorizations in the past")
	}

	// Check if the user already has on-chain funds
	hasFunds, err := h.chainClient.HasFunds(req.DesmosAddress)
	if err != nil {
		return err
	}

	if hasFunds {
		return utils.WrapErr(codes.FailedPrecondition, "you already have funds in your wallet")
	}

	// Check if the user already has an on-chain grant
	hasGrant, err := h.chainClient.HasFeeGrant(req.DesmosAddress, app.WalletAddress)
	if err != nil {
		return err
	}

	var grantTime *time.Time
	if hasGrant {
		grantTime = utils.GetTimePointer(time.Now())
	}

	// Store the request inside the database
	request := types.NewFeeGrantRequest(req.AppID, req.DesmosAddress, allowance, time.Now(), grantTime)
	return h.db.SaveFeeGrantRequest(request)
}
