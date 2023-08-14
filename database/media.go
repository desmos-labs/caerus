package database

// SaveMediaHash allows to save the given image hash inside the database
func (db *Database) SaveMediaHash(imageUrl string, hash string) error {
	stmt := `INSERT INTO images_hashes (image_url, hash) VALUES ($1, $2) ON CONFLICT (image_url) DO UPDATE SET hash = $2`
	_, err := db.SQL.Exec(stmt, imageUrl, hash)
	return err
}
