package local

import (
	"fmt"
	"os"
	"strings"
)

// ReadFile reads a file from the local file system.
// For security, we might want to restrict this in the future,
// but for a local dev tool, full access is provided.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(data), nil
}

// ListDirectory lists contents of a directory.
func ListDirectory(path string) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Contents of %s:\n\n", path))

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to list directory: %v", err)
	}

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		
		typeStr := "FILE"
		if e.IsDir() {
			typeStr = "DIR "
		}
		
		sb.WriteString(fmt.Sprintf("[%s] %s (Size: %d bytes)\n", typeStr, e.Name(), info.Size()))
	}
	
	return sb.String(), nil
}
