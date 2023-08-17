package notifications

import (
	"google.golang.org/grpc/codes"

	"github.com/desmos-labs/caerus/types"
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
		return utils.WrapErr(codes.FailedPrecondition, "application not found")
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
		return utils.NewTooManyRequestsError("notifications rate limit reached")
	}

	// Get the notification tokens for each address
	var notificationTokens []string
	for _, address := range req.UserAddresses {
		tokens, err := h.db.GetUserNotificationTokens(address)
		if err != nil {
			return err
		}

		notificationTokens = append(notificationTokens, tokens...)
	}

	// Send the notification if there are some tokens to which to send it
	if len(notificationTokens) > 0 {
		err = h.firebase.SendNotifications(app, notificationTokens, req.Notification)
		if err != nil {
			return err
		}
	}

	// Store the sent notification
	return h.db.SaveSentNotification(types.NewSentNotification(req.AppID, req.UserAddresses, req.Notification))
}
