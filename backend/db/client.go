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
	CREATE TABLE IF NOT EXISTS provider_keys_multi (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		provider_id TEXT NOT NULL,
		api_key TEXT NOT NULL,
		UNIQUE(provider_id, api_key)
	);
	CREATE TABLE IF NOT EXISTS usage_stats (
		provider_id TEXT PRIMARY KEY,
		request_count INTEGER DEFAULT 0
	);`
	
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Migrate old keys to new table
	_, _ = db.Exec(`
		INSERT OR IGNORE INTO provider_keys_multi (provider_id, api_key)
		SELECT provider_id, api_key FROM provider_keys;
	`)

	DB = db
	return nil
}

// GetKey retrieves a random key from the database for a provider. Returns empty string if not found.
func GetKey(providerID string) string {
	if DB == nil {
		return ""
	}

	var apiKey string
	// ORDER BY RANDOM() load balances the requests across available keys
	err := DB.QueryRow("SELECT api_key FROM provider_keys_multi WHERE provider_id = ? ORDER BY RANDOM() LIMIT 1", providerID).Scan(&apiKey)
	if err != nil {
		return ""
	}
	return apiKey
}

// SetKey adds a new key into the database for a provider. (Kept for backwards compatibility if needed, but renamed/repurposed to insert)
func SetKey(providerID, apiKey string) error {
	return AddKey(providerID, apiKey)
}

// AddKey adds a new key.
func AddKey(providerID, apiKey string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	upsertSQL := `
	INSERT INTO provider_keys_multi (provider_id, api_key) 
	VALUES (?, ?)
	ON CONFLICT(provider_id, api_key) DO NOTHING;
	`
	_, err := DB.Exec(upsertSQL, providerID, apiKey)
	return err
}

// DeleteKey removes a specific key.
func DeleteKey(providerID, apiKey string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := DB.Exec("DELETE FROM provider_keys_multi WHERE provider_id = ? AND api_key = ?", providerID, apiKey)
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
