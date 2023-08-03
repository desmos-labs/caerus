package database

import (
	"database/sql"
	"errors"

	"github.com/desmos-labs/caerus/types"
)

type notificationApplicationRow struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type notificationSenderRow struct {
	Token         string `db:"token"`
	ApplicationID string `db:"application_id"`
}

// GetNotificationSender returns the notification sender having the given token
func (db *Database) GetNotificationSender(token string) (*types.NotificationSender, bool, error) {
	stmt := `SELECT * FROM notification_senders WHERE token = $1`

	var senderRow notificationSenderRow
	err := db.SQL.Get(&senderRow, stmt, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	var applicationRow notificationApplicationRow
	stmt = `SELECT * FROM notification_applications WHERE id = $1`
	err = db.SQL.Get(&applicationRow, stmt, senderRow.ApplicationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	return &types.NotificationSender{
		Token: senderRow.Token,
		Application: &types.NotificationApplication{
			ID:   applicationRow.ID,
			Name: applicationRow.Name,
		},
	}, true, nil
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
