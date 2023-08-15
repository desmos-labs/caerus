package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/desmos-labs/caerus/types"
)

// --------------------------------------------------------------------------------------------------------------------
// --- NONCES
// --------------------------------------------------------------------------------------------------------------------

type nonceRow struct {
	DesmosAddress  string    `db:"desmos_address"`
	Token          string    `db:"token"`
	ExpirationTime time.Time `db:"expiration_time"`
}

// GetNonce returns all the nonce stored for the given user
func (db *Database) GetNonce(desmosAddress string, value string) (*types.EncryptedNonce, error) {
	stmt := `SELECT * FROM nonces WHERE desmos_address = $1 AND token = $2`

	var rows []nonceRow
	err := db.SQL.Select(&rows, stmt, desmosAddress, db.encryptValue(value))
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return types.NewEncryptedNonce(rows[0].DesmosAddress, rows[0].Token, rows[0].ExpirationTime), err
}

// SaveNonce stores the given nonce
func (db *Database) SaveNonce(nonce *types.Nonce) error {
	stmt := `
INSERT INTO nonces (desmos_address, token, expiration_time) 
VALUES ($1, $2, $3) ON CONFLICT (token) DO UPDATE 
    SET desmos_address = excluded.desmos_address,
        token = excluded.token,
        expiration_time = excluded.expiration_time`
	_, err := db.SQL.Exec(stmt, nonce.DesmosAddress, db.encryptValue(nonce.Value), nonce.ExpirationTime)
	return err
}

// DeleteNonce removes the given nonce from the database
func (db *Database) DeleteNonce(nonce *types.EncryptedNonce) error {
	stmt := `DELETE FROM nonces WHERE desmos_address = $1 AND token = $2`
	_, err := db.SQL.Exec(stmt, nonce.DesmosAddress, nonce.EncryptedValue)
	return err
}

// --------------------------------------------------------------------------------------------------------------------
// --- SESSIONS
// --------------------------------------------------------------------------------------------------------------------

// SaveSession stores the given session inside the database
func (db *Database) SaveSession(session *types.UserSession) error {
	stmt := `
INSERT INTO sessions (desmos_address, token, creation_time, expiration_time) 
VALUES ($1, $2, $3, $4) ON CONFLICT (token) DO UPDATE 
    SET desmos_address = excluded.desmos_address,
        token = excluded.token,
        creation_time = excluded.creation_time, 
        expiration_time = excluded.expiration_time`
	_, err := db.SQL.Exec(stmt, session.DesmosAddress, db.encryptValue(session.Token), session.CreationDate, session.ExpiryTime)
	return err
}

type sessionRow struct {
	DesmosAddress  string    `db:"desmos_address"`
	Token          string    `db:"token"`
	CreationTime   time.Time `db:"creation_time"`
	ExpirationTime time.Time `db:"expiration_time"`
}

// GetUserSession returns the session associated to the given token, if any
func (db *Database) GetUserSession(token string) (*types.EncryptedUserSession, error) {
	stmt := `SELECT * FROM sessions WHERE token = $1`

	var row sessionRow
	err := db.SQL.Get(&row, stmt, db.encryptValue(token))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return types.NewEncryptedUserSession(
		row.DesmosAddress,
		row.Token,
		row.CreationTime,
		&row.ExpirationTime,
	), nil
}

// UpdateSession updates the stored data of the given session
func (db *Database) UpdateSession(session *types.EncryptedUserSession) error {
	stmt := `UPDATE sessions SET expiration_time = $3 WHERE desmos_address = $1 AND token = $2`
	_, err := db.SQL.Exec(stmt, session.DesmosAddress, session.EncryptedToken, session.ExpiryTime)
	return err
}

// DeleteSession deletes the given session from the database
func (db *Database) DeleteSession(sessionToken string) error {
	stmt := `DELETE FROM sessions WHERE token = $1`
	_, err := db.SQL.Exec(stmt, sessionToken)
	return err
}

// DeleteAllSessions deletes all the sessions associated to the given user
func (db *Database) DeleteAllSessions(desmosAddress string) error {
	stmt := `DELETE FROM sessions WHERE desmos_address = $1`
	_, err := db.SQL.Exec(stmt, desmosAddress)
	return err
}

// --------------------------------------------------------------------------------------------------------------------
// --- USERS
// --------------------------------------------------------------------------------------------------------------------

// SaveUser stores the given user inside the database
func (db *Database) SaveUser(desmosAddress string) error {
	stmt := `INSERT INTO users (desmos_address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.SQL.Exec(stmt, desmosAddress)
	return err
}

// UpdateLoginInfo updates the login info for the user with the given address
func (db *Database) UpdateLoginInfo(desmosAddress string) error {
	stmt := `INSERT INTO users (desmos_address) VALUES ($1) ON CONFLICT (desmos_address) DO UPDATE SET last_login = NOW()`
	_, err := db.SQL.Exec(stmt, desmosAddress)
	return err
}

// DeleteUser removes the user with the given address from the database
func (db *Database) DeleteUser(desmosAddress string) error {
	stmt := `DELETE FROM users WHERE desmos_address = $1`
	_, err := db.SQL.Exec(stmt, desmosAddress)
	return err
}
