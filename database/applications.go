package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/desmos-labs/caerus/types"
)

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
	stmt := `INSERT INTO application_tokens (application_id, token_name, token_value) VALUES ($1, $2, $3)`
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
