package search

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func FetchWebsite(targetURL string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	text := string(bodyBytes)

	// Remove script and style tags
	reScript := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	text = reScript.ReplaceAllString(text, "")
	reStyle := regexp.MustCompile(`(?is)<style.*?>.*?</style>`)
	text = reStyle.ReplaceAllString(text, "")
	
	// Remove HTML tags
	reHTML := regexp.MustCompile(`(?is)<.*?>`)
	text = reHTML.ReplaceAllString(text, " ")

	// Condense whitespace
	reSpace := regexp.MustCompile(`\s+`)
	text = reSpace.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// Cap at 15000 chars to avoid blowing up context window
	if len(text) > 15000 {
		text = text[:15000] + "... [truncated]"
	}
	return text, nil
}
