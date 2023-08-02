package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

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
func (db *Database) SaveSession(session *types.Session) error {
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

// GetSession returns the session associated to the given token, if any
func (db *Database) GetSession(token string) (*types.EncryptedSession, error) {
	stmt := `SELECT * FROM sessions WHERE token = $1`

	var row sessionRow
	err := db.SQL.Get(&row, stmt, db.encryptValue(token))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return types.NewEncryptedSession(
		row.DesmosAddress,
		row.Token,
		row.CreationTime,
		&row.ExpirationTime,
	), nil
}

// UpdateSession updates the stored data of the given session
func (db *Database) UpdateSession(session *types.EncryptedSession) error {
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

// --------------------------------------------------------------------------------------------------------------------
// --- FEE GRANT REQUESTS
// --------------------------------------------------------------------------------------------------------------------

// SaveFeeGrantRequest stores a new authorization request from the given user
func (db *Database) SaveFeeGrantRequest(request types.FeeGrantRequest) error {
	stmt := `
INSERT INTO fee_grant_requests (desmos_address, request_time, grant_time) 
VALUES ($1, $2, $3)
ON CONFLICT (desmos_address) DO NOTHING `

	_, err := db.SQL.Exec(stmt, request.DesmosAddress, request.RequestTime, request.GrantTime)
	return err
}

// SetFeeGrantRequestGranted sets the fee grant requests for the given users as granted
func (db *Database) SetFeeGrantRequestGranted(users []string) error {
	stmt := `UPDATE fee_grant_requests SET grant_time = NOW() WHERE desmos_address = ANY($1)`
	_, err := db.SQL.Exec(stmt, pq.Array(users))
	return err
}

// HasFeeGrantBeenGrantedToUser returns true if the given user has already been granted the fee grant
func (db *Database) HasFeeGrantBeenGrantedToUser(user string) (bool, error) {
	stmt := `SELECT EXISTS (SELECT 1 FROM fee_grant_requests WHERE desmos_address = $1 AND grant_time IS NOT NULL)`
	var exists bool
	err := db.SQL.QueryRow(stmt, user).Scan(&exists)
	return exists, err
}

type feeGrantRequestRow struct {
	DesmosAddress string       `db:"desmos_address"`
	RequestTime   time.Time    `db:"request_time"`
	GrantTime     sql.NullTime `db:"grant_time"`
}

// GetFeeGrantRequests returns the oldest fee grant requests made from users
func (db *Database) GetFeeGrantRequests(limit int) ([]types.FeeGrantRequest, error) {
	stmt := `SELECT * FROM fee_grant_requests ORDER BY request_time LIMIT $1`

	var rows []feeGrantRequestRow
	err := db.SQL.Select(&rows, stmt, limit)
	if err != nil {
		return nil, err
	}

	var requests = make([]types.FeeGrantRequest, len(rows))
	for index, row := range rows {
		requests[index] = types.FeeGrantRequest{
			DesmosAddress: row.DesmosAddress,
			RequestTime:   row.RequestTime,
			GrantTime:     NullTimeToTime(row.GrantTime),
		}
	}

	return requests, nil
}
