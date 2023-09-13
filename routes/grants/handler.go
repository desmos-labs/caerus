package grants

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/gogoproto/proto"
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
func (h *Handler) HandleFeeGrantRequest(req *RequestFeeGrantRequest) (*RequestFeeAllowanceResponse, error) {
	// Unmarshal the allowance
	var allowance feegrant.FeeAllowanceI
	err := h.cdc.UnpackAny(req.Allowance, &allowance)
	if err != nil {
		return nil, err
	}

	// Get the app details
	app, found, err := h.db.GetApp(req.AppID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, utils.WrapErr(codes.FailedPrecondition, "application not found")
	}

	// Check if the app has granted a MsgGrantFeeAllowance permission
	hasGrantedAuthzAuthorization, err := h.chainClient.HasGrantedMsgGrantAllowanceAuthorization(app.WalletAddress)
	if err != nil {
		return nil, err
	}

	if !hasGrantedAuthzAuthorization {
		return nil, utils.WrapErr(codes.FailedPrecondition, "on-chain authorization not found")
	}

	// Check if the application has reached the number of requests
	requestRateLimit, err := h.db.GetAppFeeGrantRequestsLimit(req.AppID)
	if err != nil {
		return nil, err
	}

	requestsCount, err := h.db.GetAppFeeGrantRequestsCount(req.AppID)
	if err != nil {
		return nil, err
	}

	if requestRateLimit > 0 && requestsCount >= requestRateLimit {
		return nil, utils.NewTooManyRequestsError("number of fee grant requests allowed reached")
	}

	// Check if the user has already been granted the fee grant
	hasBeenGranted, err := h.db.HasFeeGrantBeenGrantedToUser(req.AppID, req.DesmosAddress)
	if err != nil {
		return nil, err
	}

	if hasBeenGranted {
		return nil, utils.WrapErr(codes.FailedPrecondition, "you have already been granted the authorizations in the past")
	}

	// Check if the user already has on-chain funds
	hasFunds, err := h.chainClient.HasFunds(req.DesmosAddress)
	if err != nil {
		return nil, err
	}

	if hasFunds {
		return nil, utils.WrapErr(codes.FailedPrecondition, "you already have funds in your wallet")
	}

	// Check if the user already has an on-chain grant
	hasGrant, err := h.chainClient.HasFeeGrant(req.DesmosAddress, app.WalletAddress)
	if err != nil {
		return nil, err
	}

	var grantTime *time.Time
	if hasGrant {
		grantTime = utils.GetTimePointer(time.Now())
	}

	// Store the request inside the database
	request := types.NewFeeGrantRequest(req.AppID, req.DesmosAddress, allowance, time.Now(), grantTime)
	err = h.db.SaveFeeGrantRequest(request)
	if err != nil {
		return nil, err
	}

	// Return the response
	return &RequestFeeAllowanceResponse{RequestId: request.ID}, nil
}

func (h *Handler) HandleRequestFeeGrantDetailsRequest(req *RequestFeeGrantDetailsRequest) (*GetFeeAllowanceDetailsResponse, error) {
	request, found, err := h.db.GetFeeGrantRequest(req.AppID, req.FeeGrantAllowanceRequestID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, utils.WrapErr(codes.NotFound, "fee grant request not found")
	}

	msg, ok := request.Allowance.(proto.Message)
	if !ok {
		return nil, utils.WrapErr(codes.Internal, fmt.Sprintf("cannot proto marshal %T", msg))
	}
	allowanceAny, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &GetFeeAllowanceDetailsResponse{
		UserDesmosAddress: request.DesmosAddress,
		Allowance:         allowanceAny,
		RequestTime:       request.RequestTime,
		GrantTime:         request.GrantTime,
	}, nil
}
