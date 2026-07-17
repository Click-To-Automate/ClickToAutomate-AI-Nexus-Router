package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// UsageData represents the telemetry for a single provider
type UsageData struct {
	Count       int `json:"count"`
	TokensSaved int `json:"tokens_saved"`
	TokensUsed  int `json:"tokens_used"`
}

var DB *sql.DB
var DBPath string

func InitDB(customDBPath string) error {
	var dbPath string

	if customDBPath != "" {
		dbPath = customDBPath
	} else {
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

		dbPath = filepath.Join(appDir, "router.db")
	}

	DBPath = dbPath

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
	);
	CREATE TABLE IF NOT EXISTS app_settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS semantic_cache (
		hash TEXT PRIMARY KEY,
		prompt TEXT NOT NULL,
		response TEXT NOT NULL,
		tokens_saved INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS loading_phrases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		phrase TEXT NOT NULL,
		usage_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS chat_sessions (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS chat_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS provider_latency (
		provider_id TEXT PRIMARY KEY,
		avg_latency_ms INTEGER DEFAULT 0,
		measurement_count INTEGER DEFAULT 0,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS model_tags (
		provider_id TEXT NOT NULL,
		model_id TEXT NOT NULL,
		tags TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(provider_id, model_id)
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

	// Try migrating usage_stats to include tokens_saved if it doesn't exist
	_, _ = db.Exec(`ALTER TABLE usage_stats ADD COLUMN tokens_saved INTEGER DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE usage_stats ADD COLUMN tokens_used INTEGER DEFAULT 0;`)

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

// IncrementUsage increments the request count, tokens used, and tokens saved for a given provider.
func IncrementUsage(providerID string, tokensUsed int, tokensSaved int) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	upsertSQL := `
	INSERT INTO usage_stats (provider_id, request_count, tokens_saved, tokens_used) 
	VALUES (?, 1, ?, ?)
	ON CONFLICT(provider_id) DO UPDATE SET 
		request_count = request_count + 1,
		tokens_saved = tokens_saved + ?,
		tokens_used = tokens_used + ?;
	`
	_, err := DB.Exec(upsertSQL, providerID, tokensSaved, tokensUsed, tokensSaved, tokensUsed)
	return err
}

// GetAllUsage retrieves the request counts, tokens used, and tokens saved for all providers.
func GetAllUsage() map[string]UsageData {
	usage := make(map[string]UsageData)
	if DB == nil {
		return usage
	}

	rows, err := DB.Query("SELECT provider_id, request_count, COALESCE(tokens_saved, 0), COALESCE(tokens_used, 0) FROM usage_stats")
	if err != nil {
		return usage
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var count int
		var saved int
		var used int
		if err := rows.Scan(&id, &count, &saved, &used); err == nil {
			usage[id] = UsageData{
				Count:       count,
				TokensSaved: saved,
				TokensUsed:  used,
			}
		}
	}
	return usage
}

// GetSetting retrieves a setting value by key.
func GetSetting(key string, defaultValue string) string {
	if DB == nil {
		return defaultValue
	}
	var value string
	err := DB.QueryRow("SELECT value FROM app_settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return defaultValue
	}
	return value
}

// SetSetting sets a setting value.
func SetSetting(key string, value string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := DB.Exec("INSERT INTO app_settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?", key, value, value)
	return err
}

// GetCache retrieves a cached response by hash.
func GetCache(hash string) (string, error) {
	if DB == nil {
		return "", fmt.Errorf("database not initialized")
	}
	var response string
	err := DB.QueryRow("SELECT response FROM semantic_cache WHERE hash = ?", hash).Scan(&response)
	return response, err
}

// SetCache saves a cached response.
func SetCache(hash string, prompt string, response string, tokensSaved int) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := DB.Exec("INSERT OR IGNORE INTO semantic_cache (hash, prompt, response, tokens_saved) VALUES (?, ?, ?, ?)", hash, prompt, response, tokensSaved)
	return err
}

// GetAllCache retrieves all cache entries for the UI.
func GetAllCache() []map[string]interface{} {
	var results []map[string]interface{}
	if DB == nil {
		return results
	}
	rows, err := DB.Query("SELECT hash, prompt, tokens_saved, created_at FROM semantic_cache ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		return results
	}
	defer rows.Close()

	for rows.Next() {
		var hash, prompt, createdAt string
		var tokensSaved int
		if err := rows.Scan(&hash, &prompt, &tokensSaved, &createdAt); err == nil {
			results = append(results, map[string]interface{}{
				"hash": hash,
				"prompt": prompt,
				"tokens_saved": tokensSaved,
				"created_at": createdAt,
			})
		}
	}
	return results
}

// ClearCache clears all cache entries.
func ClearCache() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := DB.Exec("DELETE FROM semantic_cache")
	return err
}

// CleanExpiredCache deletes cache entries older than the specified number of days.
func CleanExpiredCache(days int) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	query := fmt.Sprintf("DELETE FROM semantic_cache WHERE created_at < datetime('now', '-%d days')", days)
	_, err := DB.Exec(query)
	return err
}

// UpdateLatency updates the moving average latency for a provider
func UpdateLatency(providerID string, latencyMs int) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	
	upsertSQL := `
	INSERT INTO provider_latency (provider_id, avg_latency_ms, measurement_count) 
	VALUES (?, ?, 1)
	ON CONFLICT(provider_id) DO UPDATE SET 
		avg_latency_ms = ((avg_latency_ms * measurement_count) + ?) / (measurement_count + 1),
		measurement_count = measurement_count + 1,
		last_updated = CURRENT_TIMESTAMP;
	`
	_, err := DB.Exec(upsertSQL, providerID, latencyMs, latencyMs)
	return err
}

// GetLatencies gets all provider latencies
func GetLatencies() map[string]int {
	latencies := make(map[string]int)
	if DB == nil {
		return latencies
	}
	rows, err := DB.Query("SELECT provider_id, avg_latency_ms FROM provider_latency")
	if err != nil {
		return latencies
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var lat int
		if err := rows.Scan(&id, &lat); err == nil {
			latencies[id] = lat
		}
	}
	return latencies
}

// SaveChatSession saves or updates a chat session
func SaveChatSession(id string, title string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := DB.Exec("INSERT INTO chat_sessions (id, title) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET title = ?", id, title, title)
	return err
}

// SaveChatMessage saves a message to a session
func SaveChatMessage(sessionID string, role string, content string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := DB.Exec("INSERT INTO chat_messages (session_id, role, content) VALUES (?, ?, ?)", sessionID, role, content)
	return err
}

type ChatSession struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GetChatSessions gets all chat sessions
func GetChatSessions() ([]ChatSession, error) {
	var sessions []ChatSession
	if DB == nil {
		return sessions, fmt.Errorf("database not initialized")
	}
	rows, err := DB.Query("SELECT id, title, created_at FROM chat_sessions ORDER BY created_at DESC")
	if err != nil {
		return sessions, err
	}
	defer rows.Close()
	for rows.Next() {
		var s ChatSession
		if err := rows.Scan(&s.ID, &s.Title, &s.CreatedAt); err == nil {
			sessions = append(sessions, s)
		}
	}
	return sessions, nil
}

// GetChatMessages gets messages for a session
func GetChatMessages(sessionID string) ([]ChatMessage, error) {
	var msgs []ChatMessage
	if DB == nil {
		return msgs, fmt.Errorf("database not initialized")
	}
	rows, err := DB.Query("SELECT role, content FROM chat_messages WHERE session_id = ? ORDER BY id ASC", sessionID)
	if err != nil {
		return msgs, err
	}
	defer rows.Close()
	for rows.Next() {
		var m ChatMessage
		if err := rows.Scan(&m.Role, &m.Content); err == nil {
			msgs = append(msgs, m)
		}
	}
	return msgs, nil
}

// DeleteChatSession deletes a session and its messages
func DeleteChatSession(sessionID string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	// Note: We use ON DELETE CASCADE in the schema, but just in case, we'll delete both
	_, _ = DB.Exec("DELETE FROM chat_messages WHERE session_id = ?", sessionID)
	_, err := DB.Exec("DELETE FROM chat_sessions WHERE id = ?", sessionID)
	return err
}

// GetTagsForModel retrieves the tags for a specific model
func GetTagsForModel(providerID string, modelID string) ([]string, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var tagsJSON string
	err := DB.QueryRow("SELECT tags FROM model_tags WHERE provider_id = ? AND model_id = ?", providerID, modelID).Scan(&tagsJSON)
	if err != nil {
		return nil, err
	}

	var tags []string
	if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// SetTagsForModel saves tags for a specific model
func SetTagsForModel(providerID string, modelID string, tags []string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO model_tags (provider_id, model_id, tags) VALUES (?, ?, ?) ON CONFLICT(provider_id, model_id) DO UPDATE SET tags = ?", providerID, modelID, string(tagsJSON), string(tagsJSON))
	return err
}
