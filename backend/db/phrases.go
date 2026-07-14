package db

import (
	"log"
)

// InsertLoadingPhrases saves newly generated contextual phrases for a specific session.
func InsertLoadingPhrases(sessionID string, phrases []string) {
	if DB == nil {
		return
	}

	for _, p := range phrases {
		_, err := DB.Exec(`
			INSERT INTO loading_phrases (session_id, phrase, usage_count)
			VALUES (?, ?, 0)
		`, sessionID, p)
		if err != nil {
			log.Printf("Failed to insert loading phrase: %v", err)
		}
	}
}

// GetActiveLoadingPhrases retrieves phrases with usage_count < 6 for a session
// and increments their usage_count so they eventually expire.
func GetActiveLoadingPhrases(sessionID string) []string {
	if DB == nil {
		return nil
	}

	rows, err := DB.Query(`
		SELECT id, phrase FROM loading_phrases 
		WHERE session_id = ? AND usage_count < 6
		ORDER BY created_at DESC LIMIT 10
	`, sessionID)

	if err != nil {
		log.Printf("Error fetching loading phrases: %v", err)
		return nil
	}
	defer rows.Close()

	var phrases []string
	var ids []int

	for rows.Next() {
		var id int
		var phrase string
		if err := rows.Scan(&id, &phrase); err == nil {
			phrases = append(phrases, phrase)
			ids = append(ids, id)
		}
	}

	// Increment usage_count in background
	if len(ids) > 0 {
		go func(updateIDs []int) {
			for _, id := range updateIDs {
				_, _ = DB.Exec(`UPDATE loading_phrases SET usage_count = usage_count + 1 WHERE id = ?`, id)
			}
		}(ids)
	}

	return phrases
}
