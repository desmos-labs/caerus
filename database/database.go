package database

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/desmos-labs/caerus/utils"
)

type Database struct {
	url string
	SQL *sqlx.DB
}

// NewDatabase returns a new Database instance
func NewDatabase(url string) (*Database, error) {
	sqlDb, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Database{
		url: url,
		SQL: sqlDb,
	}, nil
}

// NewDatabaseFromEnvVariables returns a new Database instance reading the configuration from the environment variables
func NewDatabaseFromEnvVariables() (*Database, error) {
	databaseURI := utils.GetEnvOr(EnvDatabaseURI, "")
	if databaseURI == "" {
		return nil, fmt.Errorf("missing environment variable %s", EnvDatabaseURI)
	}

	return NewDatabase(databaseURI)
}

// encryptValue encrypts the given values so that it can be stored in the database safely
func (db *Database) encryptValue(value string) string {
	hashBz := sha256.Sum256([]byte(value))
	return base64.StdEncoding.EncodeToString(hashBz[:])
}
