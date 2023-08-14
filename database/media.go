package database

// SaveMediaHash allows to save the given image hash inside the database
func (db *Database) SaveMediaHash(imageUrl string, hash string) error {
	stmt := `INSERT INTO files_hashes (file_name, hash) VALUES ($1, $2) ON CONFLICT (file_name) DO UPDATE SET hash = $2`
	_, err := db.SQL.Exec(stmt, imageUrl, hash)
	return err
}
