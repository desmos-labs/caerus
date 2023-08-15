package notifications

import (
	"net/http"

	"github.com/desmos-labs/caerus/utils"
)

type Handler struct {
	firebase Firebase
	db       Database
}

// NewHandler allows to build a new Handler instance
func NewHandler(firebase Firebase, db Database) *Handler {
	return &Handler{
		firebase: firebase,
		db:       db,
	}
}

// HandleSendNotificationRequest handles the request to send a new notification
func (h *Handler) HandleSendNotificationRequest(req *SendAppNotificationRequest) error {
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
