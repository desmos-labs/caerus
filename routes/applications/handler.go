package applications

import (
	"github.com/desmos-labs/caerus/types"
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

// HandleRegisterAppDeviceTokenRequest handles the request to register a new application device token
func (h *Handler) HandleRegisterAppDeviceTokenRequest(req *RegisterAppDeviceTokenRequest) error {
	return h.db.SaveAppNotificationDeviceToken(types.NewAppNotificationDeviceToken(req.AppID, req.DeviceToken))
}

// HandleDeleteApplicationRequest handles the request to delete an application
func (h *Handler) HandleDeleteApplicationRequest(req *DeleteApplicationRequest) error {
	return h.db.DeleteApp(req.AppID)
}
