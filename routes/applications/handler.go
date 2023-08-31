package applications

import (
	"google.golang.org/grpc/codes"

	"github.com/desmos-labs/caerus/utils"
)

type Handler struct {
	db Database
}

// NewHandler allows to build a new Handler instance
func NewHandler(db Database) *Handler {
	return &Handler{
		db: db,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// HandleDeleteApplicationRequest handles the request to delete an application
func (h *Handler) HandleDeleteApplicationRequest(req *DeleteApplicationRequest) error {
	// Check to make sure the user can delete the app
	isAdmin, err := h.db.IsUserAdminOfApp(req.UserAddress, req.AppID)
	if err != nil {
		return err
	}

	if !isAdmin {
		return utils.WrapErr(codes.PermissionDenied, "You cannot delete this application")
	}

	return h.db.DeleteApp(req.AppID)
}
