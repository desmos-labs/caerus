package applications

import (
	"encoding/json"

	"github.com/desmos-labs/caerus/server/routes/base"
	"github.com/desmos-labs/caerus/types"
)

type Handler struct {
	*base.Handler
	db Database
}

// ParseRegisterAppDeviceTokenRequest parses the given body into a RegisterAppDeviceTokenRequest
func (h *Handler) ParseRegisterAppDeviceTokenRequest(body []byte) (*RegisterAppDeviceTokenRequest, error) {
	var req RegisterAppDeviceTokenRequest
	return &req, json.Unmarshal(body, &req)
}

// HandleRegisterAppDeviceTokenRequest handles the request to register a new application device token
func (h *Handler) HandleRegisterAppDeviceTokenRequest(req *RegisterAppDeviceTokenRequest) error {
	return h.db.SaveAppNotificationDeviceToken(types.NewAppNotificationDeviceToken(req.AppID, req.DeviceToken))
}

// --------------------------------------------------------------------------------------------------------------------

// HandleDeleteApplicationRequest handles the request to delete an application
func (h *Handler) HandleDeleteApplicationRequest(req *DeleteApplicationRequest) error {
	return h.db.DeleteApp(req.AppID)
}
