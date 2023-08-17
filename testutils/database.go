package testutils

import (
	"os"
	"path/filepath"

	"github.com/desmos-labs/desmos/v5/app"
	_ "github.com/lib/pq"

	"github.com/desmos-labs/caerus/database"
)

func CreateDatabase(pathToDbSchema string) (*database.Database, error) {
	cdc, _ := app.MakeCodecs()

	db, err := database.NewDatabase("postgres://caerus:password@localhost:6433/caerus?sslmode=disable&search_path=public", cdc)
	if err != nil {
		return nil, err
	}

	// Delete the public schema
	_, err = db.SQL.Exec(`DROP SCHEMA IF EXISTS public CASCADE;`)
	if err != nil {
		return nil, err
	}

	// Re-create the schema
	_, err = db.SQL.Exec(`CREATE SCHEMA public;`)
	if err != nil {
		return nil, err
	}

	dir, err := os.ReadDir(pathToDbSchema)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range dir {
		file, err := os.ReadFile(filepath.Join(pathToDbSchema, fileInfo.Name()))
		if err != nil {
			return nil, err
		}

		_, err = db.SQL.Exec(string(file))
		if err != nil {
			return nil, err
		}
	}

	// Create the truncate function
	stmt := `
CREATE OR REPLACE FUNCTION truncate_tables(username IN VARCHAR) RETURNS void AS $$
DECLARE
    statements CURSOR FOR
        SELECT tablename FROM pg_tables
        WHERE tableowner = username AND schemaname = 'public' 
          AND tablename != 'user_types';
BEGIN
    FOR stmt IN statements LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' CASCADE;';
    END LOOP;
END;
$$ LANGUAGE plpgsql;`
	_, err = db.SQL.Exec(stmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TruncateDatabase(db *database.Database) error {
	_, err := db.SQL.Exec(`SELECT truncate_tables('caerus')`)
	return err
}
