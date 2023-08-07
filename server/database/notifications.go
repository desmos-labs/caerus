package database

import (
	"database/sql"
	"errors"

	"github.com/desmos-labs/caerus/types"
)

// GetAppNotificationsRateLimit returns the notifications rate limit for the given application.
// If 0 is returned, it means that the application has no rate limit.
func (db *Database) GetAppNotificationsRateLimit(appID string) (uint64, error) {
	stmt := `
SELECT COALESCE(application_subscriptions.notifications_rate_limit, 0)
FROM applications
LEFT JOIN application_subscriptions ON application_subscriptions.id = applications.subscription_id
WHERE applications.id = $1;`

	var limit uint64
	err := db.SQL.Get(&limit, stmt, appID)
	return limit, err
}

// GetAppNotificationsCount returns the number of notifications that the given application
// has sent during the current day
func (db *Database) GetAppNotificationsCount(appID string) (uint64, error) {
	stmt := `
SELECT COUNT(id)
FROM notifications
WHERE application_id = $1
	  AND send_time >= CURRENT_DATE
	  AND send_time < CURRENT_DATE + INTERVAL '1 day'`

	var count uint64
	err := db.SQL.Get(&count, stmt, appID)
	return count, err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveAppNotificationDeviceToken allows to save the given device token inside the database
func (db *Database) SaveAppNotificationDeviceToken(token *types.AppNotificationDeviceToken) error {
	stmt := `
INSERT INTO application_notifications_tokens (application_id, device_token)
VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.SQL.Exec(stmt, token.AppID, token.DeviceToken)
	return err
}

// GetAppNotificationTokens returns all the device notification tokens for the application having the given ID
func (db *Database) GetAppNotificationTokens(appID string) ([]string, error) {
	stmt := `SELECT device_token FROM application_notifications_tokens WHERE application_id = $1`

	var tokens []string
	err := db.SQL.Select(&tokens, stmt, appID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return tokens, nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveUserNotificationDeviceToken allows to save the given device token inside the database
func (db *Database) SaveUserNotificationDeviceToken(token *types.UserNotificationDeviceToken) error {
	stmt := `
INSERT INTO user_notifications_tokens (user_address, device_token)
VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.SQL.Exec(stmt, token.UserAddress, token.DeviceToken)
	return err
}

// GetUserNotificationTokens returns all the device notification tokens for the user having the given address
func (db *Database) GetUserNotificationTokens(userAddress string) ([]string, error) {
	stmt := `SELECT device_token FROM user_notifications_tokens WHERE user_address = $1`

	var tokens []string
	err := db.SQL.Select(&tokens, stmt, userAddress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return tokens, nil
}
