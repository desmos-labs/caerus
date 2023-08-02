package notifications

import (
	"encoding/json"
	"net/http"

	"github.com/desmos-labs/caerus/server/routes/base"
	"github.com/desmos-labs/caerus/server/utils"
	"github.com/desmos-labs/caerus/types"
)

type Handler struct {
	*base.Handler
	cfg      *Config
	firebase Firebase
	db       Database
}

func NewHandler(cfg *Config, firebase Firebase, db Database) *Handler {
	return &Handler{
		Handler:  base.NewHandler(db),
		cfg:      cfg,
		firebase: firebase,
		db:       db,
	}
}

// NewHandlerFromEnvVariables builds a new Handler instance reading the configuration from the env variables
func NewHandlerFromEnvVariables(firebase Firebase, db Database) *Handler {
	return NewHandler(ReadConfigFromEnvVariables(), firebase, db)
}

// --------------------------------------------------------------------------------------------------------------------

// ParseRegisterDeviceTokenRequest parses the given body into a RegisterDeviceTokenRequest
func (h *Handler) ParseRegisterDeviceTokenRequest(body []byte) (*RegisterDeviceTokenRequest, error) {
	var req RegisterDeviceTokenRequest
	err := json.Unmarshal(body, &req)
	return &req, err
}

// HandleRegisterDeviceTokenRequest handles the request to register a new device token
func (h *Handler) HandleRegisterDeviceTokenRequest(req *RegisterDeviceTokenRequest) error {
	return h.db.SaveNotificationDeviceToken(types.NewNotificationDeviceToken(req.UserAddress, req.DeviceToken))
}

// --------------------------------------------------------------------------------------------------------------------

// ParseSendNotificationRequest parses the given request body into a SendNotificationRequest
func (h *Handler) ParseSendNotificationRequest(body []byte) (*SendNotificationRequest, error) {
	var req SendNotificationRequest
	err := json.Unmarshal(body, &req)
	return &req, err
}

// HandleSendNotificationRequest handles the request to send a new notification
func (h *Handler) HandleSendNotificationRequest(req *SendNotificationRequest) error {
	if h.cfg.RequiresAuthentication {
		// Make sure the user can send notifications
		canSend, err := h.db.CanSendNotifications(req.Token)
		if err != nil {
			return err
		}

		if !canSend {
			return utils.WrapErr(http.StatusUnauthorized, "you cannot send notifications")
		}
	}

	// Send the notification
	return h.firebase.SendNotifications(req.DeviceTokens, req.Notification)
}
