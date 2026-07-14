package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level    string `json:"level"`
	Message  string `json:"message"`
}

type MemoryLogBuffer struct {
	mu      sync.Mutex
	entries []LogEntry
	maxSize int
}

func NewMemoryLogBuffer(maxSize int) *MemoryLogBuffer {
	return &MemoryLogBuffer{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

func (m *MemoryLogBuffer) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	line := string(p)
	line = strings.TrimSpace(line)
	if line == "" {
		return len(p), nil
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "INFO",
		Message:   line,
	}

	if len(m.entries) >= m.maxSize {
		m.entries = m.entries[1:]
	}
	m.entries = append(m.entries, entry)

	return len(p), nil
}

func (m *MemoryLogBuffer) GetLogs() []LogEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	res := make([]LogEntry, len(m.entries))
	copy(res, m.entries)
	return res
}

var GlobalLogBuffer = NewMemoryLogBuffer(200)

// HandleLogs returns the real log entries from memory buffer.
func HandleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logs := GlobalLogBuffer.GetLogs()
	if logs == nil {
		logs = []LogEntry{}
	}

	json.NewEncoder(w).Encode(logs)
}