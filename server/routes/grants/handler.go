package grants

import (
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

// HandleFeeGrantRequest handles the request of a fee grant from the user having the given token
func (h *Handler) HandleFeeGrantRequest(token string) error {
	session, err := h.GetSession(token)
	if err != nil {
		return err
	}

	// Check if the user has already been granted the fee grant
	hasBeenGranted, err := h.db.HasFeeGrantBeenGrantedToUser(session.DesmosAddress)
	if err != nil {
		return err
	}

	if hasBeenGranted {
		return serverutils.WrapErr(http.StatusBadRequest, "You have already been granted the authorizations in the past")
	}

	// Check if the user already has on-chain funds
	hasFunds, err := h.client.HasFunds(session.DesmosAddress)
	if err != nil {
		return err
	}

	if hasFunds {
		return serverutils.WrapErr(http.StatusBadRequest, "You already have funds in your wallet")
	}

	// Check if the user already has an on-chain grant
	hasGrant, err := h.client.HasFeeGrant(session.DesmosAddress)
	if err != nil {
		return err
	}

	var grantTime *time.Time
	if hasGrant {
		grantTime = utils.GetTimePointer(time.Now())
	}

	// Store the request inside the database
	return h.db.SaveFeeGrantRequest(types.NewFeeGrantRequest(session.DesmosAddress, time.Now(), grantTime))
}
