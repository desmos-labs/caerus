package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/gogoproto/proto"
	"github.com/lib/pq"

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
INSERT INTO fee_grant_requests (id, application_id, allowance, grantee_address, request_time, grant_time) 
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT ON CONSTRAINT unique_application_user_fee_grant_request DO NOTHING `

	protoAllowance, ok := request.Allowance.(proto.Message)
	if !ok {
		return fmt.Errorf("cannot proto marshal %T", request.Allowance)
	}

	allowanceBz, err := db.cdc.MarshalInterfaceJSON(protoAllowance)
	if err != nil {
		return fmt.Errorf("cannot marshal %T to json", protoAllowance)
	}

	_, err = db.SQL.Exec(stmt,
		request.ID,
		request.AppID,
		string(allowanceBz),
		request.DesmosAddress,
		request.RequestTime,
		TimeToNullTime(request.GrantTime),
	)
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

// SetFeeGrantRequestsGranted sets the fee grant requests having the given ids as granted
func (db *Database) SetFeeGrantRequestsGranted(ids []string) error {
	stmt := `UPDATE fee_grant_requests SET grant_time = NOW() WHERE id = ANY($1::TEXT[])`
	_, err := db.SQL.Exec(stmt, pq.StringArray(ids))
	return err
}

// HasFeeGrantBeenGrantedToUser returns true if the given user has already been granted the fee grant
func (db *Database) HasFeeGrantBeenGrantedToUser(appID string, user string) (bool, error) {
	stmt := `
SELECT EXISTS (
SELECT 1 
FROM fee_grant_requests 
WHERE application_id = $1
  AND grantee_address = $2
  AND grant_time IS NOT NULL
)`
	var exists bool
	err := db.SQL.Get(&exists, stmt, appID, user)
	return exists, err
}

type feeGrantRequestRow struct {
	ID             string       `db:"id"`
	AppID          string       `db:"application_id"`
	GranteeAddress string       `db:"grantee_address"`
	Allowance      []byte       `db:"allowance"`
	RequestTime    time.Time    `db:"request_time"`
	GrantTime      sql.NullTime `db:"grant_time"`
}

// GetNotGrantedFeeGrantRequests returns the oldest fee grant requests made from users
func (db *Database) GetNotGrantedFeeGrantRequests(limit int) ([]types.FeeGrantRequest, error) {
	stmt := `SELECT * FROM fee_grant_requests WHERE grant_time IS NULL ORDER BY request_time LIMIT $1`

	var rows []feeGrantRequestRow
	err := db.SQL.Select(&rows, stmt, limit)
	if err != nil {
		return nil, err
	}

	var requests = make([]types.FeeGrantRequest, len(rows))
	for index, row := range rows {
		var allowance feegrant.FeeAllowanceI
		err := db.cdc.UnmarshalInterfaceJSON(row.Allowance, &allowance)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal allowance: %s", err)
		}

		requests[index] = types.FeeGrantRequest{
			ID:            row.ID,
			AppID:         row.AppID,
			DesmosAddress: row.GranteeAddress,
			Allowance:     allowance,
			RequestTime:   row.RequestTime,
			GrantTime:     NullTimeToTime(row.GrantTime),
		}
	}

	return requests, nil
}

// GetFeeGrantRequest returns the fee grant request having the given id and application id
func (db *Database) GetFeeGrantRequest(appID string, requestID string) (*types.FeeGrantRequest, bool, error) {
	stmt := `SELECT * FROM fee_grant_requests WHERE application_id = $1 AND id = $2`

	var rows []feeGrantRequestRow
	err := db.SQL.Select(&rows, stmt, appID, requestID)
	if err != nil {
		return nil, false, err
	}

	if len(rows) == 0 {
		return nil, false, nil
	}

	row := rows[0]
	var allowance feegrant.FeeAllowanceI
	err = db.cdc.UnmarshalInterfaceJSON(row.Allowance, &allowance)
	if err != nil {
		return nil, false, fmt.Errorf("cannot unmarshal allowance: %s", err)
	}

	return &types.FeeGrantRequest{
		ID:            row.ID,
		AppID:         row.AppID,
		DesmosAddress: row.GranteeAddress,
		Allowance:     allowance,
		RequestTime:   row.RequestTime,
		GrantTime:     NullTimeToTime(row.GrantTime),
	}, true, nil
}
