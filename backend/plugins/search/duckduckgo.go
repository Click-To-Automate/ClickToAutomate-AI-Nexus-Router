package search

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	utls "github.com/refraction-networking/utls"
)

type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type SERPResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}

// Creates an HTTP client that perfectly fakes a Google Chrome TLS fingerprint
func newUTLSClient() *http.Client {
	transport := &http.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := net.Dialer{Timeout: 10 * time.Second}
			rawConn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}

			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				host = addr
			}

			uConn := utls.UClient(rawConn, &utls.Config{
				ServerName: host,
			}, utls.HelloChrome_Auto)
			
			err = uConn.BuildHandshakeState()
			if err != nil {
				rawConn.Close()
				return nil, err
			}
			for _, ext := range uConn.Extensions {
				if alpn, ok := ext.(*utls.ALPNExtension); ok {
					alpn.AlpnProtocols = []string{"http/1.1"}
					break
				}
			}

			err = uConn.Handshake()
			if err != nil {
				rawConn.Close()
				return nil, err
			}

			return uConn, nil
		},
		ForceAttemptHTTP2: false,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}
}

func ScrapeSERP(query string) (*SERPResponse, error) {
	result := SERPResponse{
		Query:   query,
		Results: make([]SearchResult, 0),
	}

	client := newUTLSClient()
	data := url.Values{}
	data.Set("q", query)

	req, err := http.NewRequest("POST", "https://html.duckduckgo.com/html/", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Origin", "https://html.duckduckgo.com")
	req.Header.Add("Referer", "https://html.duckduckgo.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "Unfortunately, bots use DuckDuckGo too") {
		return nil, fmt.Errorf("bot protection triggered despite uTLS spoofing")
	}

	blocks := strings.Split(bodyString, "<div class=\"result ")
	if len(blocks) > 0 {
		blocks = blocks[1:]
	}

	titleRegex := regexp.MustCompile(`(?is)<a[^>]*class="result__a"[^>]*>(.*?)</a>`)
	urlRegex := regexp.MustCompile(`(?is)<a[^>]*class="result__url"[^>]*href="([^"]+)"`)
	snippetRegex := regexp.MustCompile(`(?is)<a[^>]*class="result__snippet"[^>]*>(.*?)</a>`)
	tagStripper := regexp.MustCompile(`(?i)<[^>]*>`)

	for _, block := range blocks {
		var item SearchResult

		if match := titleRegex.FindStringSubmatch(block); len(match) > 1 {
			item.Title = strings.TrimSpace(tagStripper.ReplaceAllString(match[1], ""))
		}
		
		if match := urlRegex.FindStringSubmatch(block); len(match) > 1 {
			rawURL := match[1]
			if strings.Contains(rawURL, "uddg=") {
				parts := strings.Split(rawURL, "uddg=")
				if len(parts) > 1 {
					decoded, _ := url.QueryUnescape(strings.Split(parts[1], "&")[0])
					item.URL = decoded
				}
			} else {
				item.URL = rawURL
			}
		}

		if match := snippetRegex.FindStringSubmatch(block); len(match) > 1 {
			item.Snippet = strings.TrimSpace(tagStripper.ReplaceAllString(match[1], ""))
		}

		if item.URL != "" && item.Title != "" {
			result.Results = append(result.Results, item)
		}
	}

	return &result, nil
}
