package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/desmos-labs/caerus/types"
)

// SaveAppSubscription allows to save the given application subscription inside the database.
func (db *Database) SaveAppSubscription(subscription types.ApplicationSubscription) error {
	stmt := `
INSERT INTO application_subscriptions (id, subscription_name, fee_grant_rate_limit, notifications_rate_limit) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE 
    SET subscription_name = excluded.subscription_name, 
        fee_grant_rate_limit = excluded.fee_grant_rate_limit, 
        notifications_rate_limit = excluded.notifications_rate_limit`
	_, err := db.SQL.Exec(stmt, subscription.ID, subscription.Name, subscription.FeeGrantLimit, subscription.NotificationsRateLimit)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveApp allows to save the given application inside the database.
func (db *Database) SaveApp(app types.Application) error {
	// Save the application
	stmt := `
INSERT INTO applications (id, name, wallet_address, subscription_id, creation_time) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE
	SET name = excluded.name,
	    wallet_address = excluded.wallet_address,
	    subscription_id = excluded.subscription_id,
	    creation_time = excluded.creation_time`
	_, err := db.SQL.Exec(stmt, app.ID, app.Name, app.WalletAddress, app.SubscriptionID, app.CreationTime)
	if err != nil {
		return err
	}

	// Save the admins
	stmt = `INSERT INTO application_admins (application_id, user_address) VALUES `

	var args []interface{}
	for i, admin := range app.Admins {
		ai := i * 2

		stmt += fmt.Sprintf("($%d, $%d),", ai+1, ai+2)
		args = append(args, app.ID, admin)
	}
	stmt = stmt[:len(stmt)-1] + ` ON CONFLICT DO NOTHING`
	_, err = db.SQL.Exec(stmt, args...)
	return err
}

// SetAppAdmin sets the given user as admin of the application having the given id
func (db *Database) SetAppAdmin(appID string, userAddress string) error {
	stmt := `INSERT INTO application_admins (application_id, user_address) VALUES ($1, $2)`
	_, err := db.SQL.Exec(stmt, appID, userAddress)
	return err
}

type applicationRow struct {
	ID             string         `db:"id"`
	Name           string         `db:"name"`
	WalletAddress  sql.NullString `db:"wallet_address"`
	SubscriptionID sql.NullInt64  `db:"subscription_id"`
	CreationTime   time.Time      `db:"creation_time"`
}

// GetApp returns the application having the given id, if any.
// If no application having the given id is found, returns false as second value.
func (db *Database) GetApp(appID string) (*types.Application, bool, error) {
	stmt := `SELECT * FROM applications WHERE id = $1`

	var row applicationRow
	err := db.SQL.Get(&row, stmt, appID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var admins []string
	stmt = `SELECT user_address FROM application_admins WHERE application_id = $1`
	err = db.SQL.Select(&admins, stmt, appID)
	if err != nil {
		return nil, false, err
	}

	return types.NewApplication(
		row.ID,
		row.Name,
		NullStringToString(row.WalletAddress),
		NullIntToUint64(row.SubscriptionID),
		admins,
		row.CreationTime,
	), true, nil
}

// IsUserAdminOfApp tells whether the given user is an administrator of the app having the given id
func (db *Database) IsUserAdminOfApp(userAddress string, appID string) (bool, error) {
	stmt := `SELECT EXISTS (SELECT 1 FROM application_admins WHERE user_address = $1 AND application_id = $2)`
	var exists bool
	err := db.SQL.Get(&exists, stmt, userAddress, appID)
	return exists, err
}

// DeleteApp deletes the application having the given id, if any.
func (db *Database) DeleteApp(appID string) error {
	stmt := `DELETE FROM applications WHERE id = $1`
	_, err := db.SQL.Exec(stmt, appID)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveAppToken allows to save the given application token inside the database
func (db *Database) SaveAppToken(token types.AppToken) error {
	stmt := `
INSERT INTO application_tokens (application_id, token_name, token_value)
VALUES ($1, $2, $3) 
ON CONFLICT ON CONSTRAINT unique_token_name_for_application DO UPDATE 
    SET token_value = excluded.token_value
`
	_, err := db.SQL.Exec(stmt, token.AppID, token.Name, db.encryptValue(token.Value))
	return err
}

type appTokenRow struct {
	ID           uint64    `db:"id"`
	AppID        string    `db:"application_id"`
	TokenName    string    `db:"token_name"`
	TokenValue   string    `db:"token_value"`
	CreationTime time.Time `db:"creation_time"`
}

// GetAppToken returns the application token having the given value, if any.
// The value is encrypted before being used to query the database.
func (db *Database) GetAppToken(token string) (*types.EncryptedAppToken, error) {
	stmt := `SELECT * FROM application_tokens WHERE token_value = $1`

	var row appTokenRow
	err := db.SQL.Get(&row, stmt, db.encryptValue(token))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return types.NewEncryptedAppToken(
		row.AppID,
		row.TokenName,
		row.TokenValue,
		row.CreationTime,
	), nil
}
