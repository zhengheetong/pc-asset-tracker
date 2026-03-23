package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB creates the local database file and table if they don't exist
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./local_cache.db")
	if err != nil {
		return nil, err
	}

	// The schema remains the same as we only need to store the hash for comparison
	query := `
	CREATE TABLE IF NOT EXISTS pc_history (
		serial TEXT PRIMARY KEY,
		hardware_hash TEXT,
		last_upload DATETIME
	);`

	_, err = db.Exec(query)
	return db, err
}

// GenerateHash creates a unique SHA-256 fingerprint of the current hardware specs
func GenerateHash(specs PCSpecs) string {
	// Added 3 more %s for the Tag1, Tag2, and Tag3 fields
	combined := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s",
		specs.CPU,
		specs.RAMTotal,
		specs.RAMModules,
		specs.Disks,
		specs.Serial,
		specs.Tag1,
		specs.Tag2,
		specs.Tag3,
	)

	hash := sha256.Sum256([]byte(combined))
	return fmt.Sprintf("%x", hash)
}

// HasHardwareChanged checks if the current fingerprint matches the one in the DB
func HasHardwareChanged(db *sql.DB, serial string, currentHash string) bool {
	var lastHash string
	err := db.QueryRow("SELECT hardware_hash FROM pc_history WHERE serial = ?", serial).Scan(&lastHash)

	if err == sql.ErrNoRows {
		// New entry required
		return true
	}

	return lastHash != currentHash
}

// UpdateLocalHash saves the new fingerprint after a successful upload
func UpdateLocalHash(db *sql.DB, serial string, newHash string) error {
	query := `
	INSERT INTO pc_history (serial, hardware_hash, last_upload) 
	VALUES (?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(serial) DO UPDATE SET 
		hardware_hash = excluded.hardware_hash,
		last_upload = CURRENT_TIMESTAMP;`

	_, err := db.Exec(query, serial, newHash)
	return err
}
