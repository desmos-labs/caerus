package grants

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/desmos-labs/caerus/server/chain"
	"github.com/desmos-labs/caerus/server/routes/base"
	serverutils "github.com/desmos-labs/caerus/server/utils"
	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

type Handler struct {
	base.Handler
	client *chain.Client
	db     Database
}

// NewHandler returns a new Handler instance
func NewHandler(client *chain.Client, db Database) *Handler {
	return &Handler{
		client: client,
		db:     db,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// ParseRequestFeeGrantRequest parses the given request body into a RequestFeeGrantRequest instance
func (h *Handler) ParseRequestFeeGrantRequest(body []byte) (*RequestFeeGrantRequest, error) {
	var req RequestFeeGrantRequest
	return &req, json.Unmarshal(body, &req)
}

// HandleFeeGrantRequest handles the request of a fee grant from the user having the given token
func (h *Handler) HandleFeeGrantRequest(req *RequestFeeGrantRequest) error {
	// Get the app details
	app, found, err := h.db.GetApp(req.AppID)
	if err != nil {
		return err
	}

	if !found {
		return serverutils.WrapErr(http.StatusNotFound, "Application not found")
	}

	// Check if the app has granted a MsgGrantFeeAllowance permission
	hasGrantedAuthorization, err := h.client.HasGrantedMsgGrantAllowanceAuthorization(app.WalletAddress)
	if err != nil {
		return err
	}

	if !hasGrantedAuthorization {
		return serverutils.WrapErr(http.StatusBadRequest, "On-chain authorization not found")
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
		return serverutils.NewTooManyRequestsError("Number of fee grant requests allowed reached")
	}

	// Check if the user has already been granted the fee grant
	hasBeenGranted, err := h.db.HasFeeGrantBeenGrantedToUser(req.AppID, req.DesmosAddress)
	if err != nil {
		return err
	}

	if hasBeenGranted {
		return serverutils.WrapErr(http.StatusBadRequest, "You have already been granted the authorizations in the past")
	}

	// Check if the user already has on-chain funds
	hasFunds, err := h.client.HasFunds(req.DesmosAddress)
	if err != nil {
		return err
	}

	if hasFunds {
		return serverutils.WrapErr(http.StatusBadRequest, "You already have funds in your wallet")
	}

	// Check if the user already has an on-chain grant
	hasGrant, err := h.client.HasFeeGrant(req.DesmosAddress, app.WalletAddress)
	if err != nil {
		return err
	}

	var grantTime *time.Time
	if hasGrant {
		grantTime = utils.GetTimePointer(time.Now())
	}

	// Store the request inside the database
	return h.db.SaveFeeGrantRequest(types.NewFeeGrantRequest(req.AppID, req.DesmosAddress, time.Now(), grantTime))
}
