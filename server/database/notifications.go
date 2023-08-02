package database

import (
	"database/sql"
	"errors"

	"github.com/desmos-labs/caerus/types"
)

// CanSendNotifications tells whether the user that is authenticating using the given token
// can send notifications or not
func (db *Database) CanSendNotifications(token string) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM notification_senders WHERE token = $1)`
	var exists bool
	err := db.SQL.Get(&exists, stmt, token)
	return exists, err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveNotificationDeviceToken allows to save the given device token inside the database
func (db *Database) SaveNotificationDeviceToken(token *types.NotificationDeviceToken) error {
	stmt := `INSERT INTO notifications_tokens (user_address, device_token) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.SQL.Exec(stmt, token.UserAddress, token.DeviceToken)
	return err
}

// GetUserNotificationTokens returns all the device notification tokens for the user having the given address
func (db *Database) GetUserNotificationTokens(userAddress string) ([]string, error) {
	stmt := `SELECT device_token FROM notifications_tokens WHERE user_address = $1`

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
