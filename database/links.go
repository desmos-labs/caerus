package database

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/desmos-labs/caerus/types"
)

// SaveCreatedDeepLink allows to save the given link inside the database
func (db *Database) SaveCreatedDeepLink(link *types.CreatedDeepLink) error {
	stmt := `
INSERT INTO deep_links (id, application_id, link_url, link_config, creation_time)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING`

	configBz, err := json.Marshal(link.Config)
	if err != nil {
		return err
	}

	_, err = db.SQL.Exec(stmt,
		link.ID,
		link.AppID,
		link.URL,
		string(configBz),
		link.CreationTime,
	)
	return err
}

// GetDeepLinkConfig allows to return the types.LinkConfig associated to the link having the given URL, if any
func (db *Database) GetDeepLinkConfig(url string) (*types.LinkConfig, error) {
	stmt := `SELECT link_config FROM deep_links WHERE link_url = $1`

	var configBz []byte
	err := db.SQL.Get(&configBz, stmt, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var config types.LinkConfig
	err = json.Unmarshal(configBz, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetAppDeepLinksRateLimit returns the deep links rate limit for the given application.
// If 0 is returned, it means that the application has no rate limit.
func (db *Database) GetAppDeepLinksRateLimit(appID string) (uint64, error) {
	stmt := `
SELECT COALESCE(application_subscriptions.deep_links_rate_limit, 0)
FROM applications
LEFT JOIN application_subscriptions ON application_subscriptions.id = applications.subscription_id
WHERE applications.id = $1;`

	var limit uint64
	err := db.SQL.Get(&limit, stmt, appID)
	return limit, err
}

// GetAppDeepLinksCount returns the number of deep links that the given application
// has created during the current day
func (db *Database) GetAppDeepLinksCount(appID string) (uint64, error) {
	stmt := `
SELECT COUNT(id)
FROM deep_links
WHERE application_id = $1
	  AND creation_time >= CURRENT_DATE
	  AND creation_time < CURRENT_DATE + INTERVAL '1 day'`

	var count uint64
	err := db.SQL.Get(&count, stmt, appID)
	return count, err
}
