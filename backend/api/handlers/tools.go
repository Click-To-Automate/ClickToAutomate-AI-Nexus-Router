package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ainexusrouter-core/plugins/search"
)

type ExecuteToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ExecuteToolResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func HandleExecuteTool(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req ExecuteToolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var result string
	var err error

	if req.Name == "search_duckduckgo" {
		result, err = executeDuckDuckGoScraper(req.Arguments)
	} else if req.Name == "fetch_website" {
		url, ok := req.Arguments["url"].(string)
		if !ok {
			err = fmt.Errorf("missing or invalid url argument")
		} else {
			result, err = executeFetchWebsite(url)
		}
	} else {
		err = fmt.Errorf("unknown tool: %s", req.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(ExecuteToolResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(ExecuteToolResponse{Result: result})
}

func executeDuckDuckGoScraper(args map[string]interface{}) (string, error) {
	query, ok := args["searchQuery"].(string)
	if !ok {
		return "", fmt.Errorf("missing searchQuery argument")
	}

	result, err := search.ScrapeSERP(query)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func executeFetchWebsite(targetURL string) (string, error) {
	return search.FetchWebsite(targetURL)
}
