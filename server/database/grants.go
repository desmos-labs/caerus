package database

import (
	"database/sql"
	"time"

	"github.com/desmos-labs/caerus/types"
)

// GetAppFeeGrantRequestsLimit returns the maximum number of fee grant requests
// that the given application can perform in a day.
// If 0 is returned, it means that the application has no limit.
func (db *Database) GetAppFeeGrantRequestsLimit(appID string) (uint64, error) {
	stmt := `
SELECT COALESCE(application_subscriptions.fee_grant_rate_limit, 0)
FROM applications
LEFT JOIN application_subscriptions ON application_subscriptions.id = applications.subscription_id
WHERE applications.id = $1;`

	var limit uint64
	err := db.SQL.Get(&limit, stmt, appID)
	return limit, err
}

// GetAppFeeGrantRequestsCount returns the number of fee grant requests
// that the given application has performed during the current day
func (db *Database) GetAppFeeGrantRequestsCount(appID string) (uint64, error) {
	stmt := `
SELECT COUNT(id) 
FROM fee_grant_requests 
WHERE application_id = $1 
  AND request_time >= CURRENT_DATE 
  AND request_time < CURRENT_DATE + INTERVAL '1 day'`

	var count uint64
	err := db.SQL.Get(&count, stmt, appID)
	return count, err
}

// -------------------------------------------------------------------------------------------------------------------

// SaveFeeGrantRequest stores a new authorization request from the given user
func (db *Database) SaveFeeGrantRequest(request types.FeeGrantRequest) error {
	stmt := `
INSERT INTO fee_grant_requests (application_id, grantee_address, grant_time) 
VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT unique_application_user_fee_grant_request DO NOTHING `

	_, err := db.SQL.Exec(stmt, request.AppID, request.GrantTime, request.GrantTime)
	return err
}

// SetFeeGrantRequestGranted sets the fee grant requests for the given users as granted
func (db *Database) SetFeeGrantRequestGranted(appID string, userAddress string) error {
	stmt := `
UPDATE fee_grant_requests 
SET grant_time = NOW()
WHERE application_id = $1
  AND grantee_address = $2`
	_, err := db.SQL.Exec(stmt, appID, userAddress)
	return err
}

// HasFeeGrantBeenGrantedToUser returns true if the given user has already been granted the fee grant
func (db *Database) HasFeeGrantBeenGrantedToUser(appID string, user string) (bool, error) {
	stmt := `
SELECT EXISTS (
SELECT 1 
FROM fee_grant_requests 
WHERE application_id = $1
  AND grantee_address = $1
  AND grant_time IS NOT NULL
)`
	var exists bool
	err := db.SQL.Get(&exists, stmt, appID, user)
	return exists, err
}

type feeGrantRequestRow struct {
	ID             uint64       `db:"id"`
	AppID          string       `db:"application_id"`
	GranteeAddress string       `db:"grantee_address"`
	RequestTime    time.Time    `db:"request_time"`
	GrantTime      sql.NullTime `db:"grant_time"`
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
		requests[index] = types.NewFeeGrantRequest(
			row.AppID,
			row.GranteeAddress,
			row.RequestTime,
			NullTimeToTime(row.GrantTime),
		)
	}

	return requests, nil
}