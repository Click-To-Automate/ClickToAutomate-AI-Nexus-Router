package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	// 1. Resolve path to C:/Users/<username>/Documents/.cta-ai-nexus
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %v", err)
	}

	appDir := filepath.Join(homeDir, "Documents", ".cta-ai-nexus")
	
	// Ensure directory exists
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create app directory: %v", err)
	}

	dbPath := filepath.Join(appDir, "router.db")

	// 2. Open SQLite database using modernc.org/sqlite
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// 3. Auto-migrate schema
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS provider_keys (
		provider_id TEXT PRIMARY KEY,
		api_key TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS usage_stats (
		provider_id TEXT PRIMARY KEY,
		request_count INTEGER DEFAULT 0
	);`
	
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	DB = db
	return nil
}

// GetKey retrieves a key from the database. Returns empty string if not found.
func GetKey(providerID string) string {
	if DB == nil {
		return ""
	}

	var apiKey string
	err := DB.QueryRow("SELECT api_key FROM provider_keys WHERE provider_id = ?", providerID).Scan(&apiKey)
	if err != nil {
		return "" // Returns empty string on sql.ErrNoRows or other errors
	}
	return apiKey
}

// SetKey upserts a key into the database.
func SetKey(providerID, apiKey string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	upsertSQL := `
	INSERT INTO provider_keys (provider_id, api_key) 
	VALUES (?, ?)
	ON CONFLICT(provider_id) DO UPDATE SET api_key=excluded.api_key;
	`
	_, err := DB.Exec(upsertSQL, providerID, apiKey)
	return err
}

// IncrementUsage increments the request count for a given provider.
func IncrementUsage(providerID string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	upsertSQL := `
	INSERT INTO usage_stats (provider_id, request_count) 
	VALUES (?, 1)
	ON CONFLICT(provider_id) DO UPDATE SET request_count = request_count + 1;
	`
	_, err := DB.Exec(upsertSQL, providerID)
	return err
}

// GetAllUsage retrieves the request counts for all providers.
func GetAllUsage() map[string]int {
	usage := make(map[string]int)
	if DB == nil {
		return usage
	}

	rows, err := DB.Query("SELECT provider_id, request_count FROM usage_stats")
	if err != nil {
		return usage
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var count int
		if err := rows.Scan(&id, &count); err == nil {
			usage[id] = count
		}
	}
	return usage
}
