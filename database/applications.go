package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/desmos-labs/caerus/types"
)

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

	return types.NewApplication(
		row.ID,
		row.Name,
		NullStringToString(row.WalletAddress),
		NullIntToUint64(row.SubscriptionID),
		row.CreationTime,
	), true, nil
}

// CanDeleteApp tells whether the given user can delete the application having the given id.
func (db *Database) CanDeleteApp(userAddress string, appID string) (bool, error) {
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
