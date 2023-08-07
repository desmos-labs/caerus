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
	firebase Firebase
	db       Database
}

// NewHandler allows to build a new Handler instance
func NewHandler(firebase Firebase, db Database) *Handler {
	return &Handler{
		Handler:  base.NewHandler(db),
		firebase: firebase,
		db:       db,
	}
}

// --------------------------------------------------------------------------------------------------------------------

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

// ParseRegisterUserDeviceTokenRequest parses the given body into a RegisterUserDeviceTokenRequest
func (h *Handler) ParseRegisterUserDeviceTokenRequest(body []byte) (*RegisterUserDeviceTokenRequest, error) {
	var req RegisterUserDeviceTokenRequest
	err := json.Unmarshal(body, &req)
	return &req, err
}

// HandleRegisterUserDeviceTokenRequest handles the request to register a new device token
func (h *Handler) HandleRegisterUserDeviceTokenRequest(req *RegisterUserDeviceTokenRequest) error {
	return h.db.SaveUserNotificationDeviceToken(types.NewUserNotificationDeviceToken(req.UserAddress, req.DeviceToken))
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
	// Get the application details
	app, found, err := h.db.GetApp(req.AppID)
	if err != nil {
		return err
	}

	if !found {
		return utils.WrapErr(http.StatusNotFound, "Application not found")
	}

	// Make sure the app has not reached the rate limit
	notificationsRateLimit, err := h.db.GetAppNotificationsRateLimit(req.AppID)
	if err != nil {
		return err
	}

	notificationsCount, err := h.db.GetAppNotificationsCount(req.AppID)
	if err != nil {
		return err
	}

	if notificationsRateLimit > 0 && notificationsCount >= notificationsRateLimit {
		return utils.NewTooManyRequestsError("Notifications rate limit reached")
	}

	// Send the notification
	return h.firebase.SendNotifications(app, req.DeviceTokens, req.Notification)
}
