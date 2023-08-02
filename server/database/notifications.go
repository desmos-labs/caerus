package database

import (
	"database/sql"
	"errors"

	"github.com/desmos-labs/caerus/types"
)

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
