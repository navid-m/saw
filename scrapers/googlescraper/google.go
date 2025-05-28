package googlescraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"swarm/models"
	"swarm/models/origins"
)

func ScrapeGoogle(query, apiKey, searchEngineID string) ([]models.SearchResult, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Google API key is required")
	}

	if searchEngineID == "" {
		return nil, fmt.Errorf("Google Custom Search Engine ID is required")
	}

	resp, shouldReturn, result, err := retrieveResults(apiKey, searchEngineID, query)
	if shouldReturn {
		return result, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status: %d %s", resp.StatusCode, resp.Status)
	}

	var googleResp GoogleSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	var results []models.SearchResult
	for _, item := range googleResp.Items {
		results = append(results, models.SearchResult{
			Origin:      origins.Google,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Snippet,
		})
	}

	return results, nil
}

func retrieveResults(apiKey string, searchEngineID string, query string) (*http.Response, bool, []models.SearchResult, error) {
	var (
		baseURL = "https://www.googleapis.com/customsearch/v1"
		params  = url.Values{}
	)

	params.Add("key", apiKey)
	params.Add("cx", searchEngineID)
	params.Add("q", query)
	params.Add("num", "10")

	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(apiURL)

	if err != nil {
		return nil, true, nil, fmt.Errorf("failed to make API request: %w", err)
	}

	defer resp.Body.Close()
	return resp, false, nil, nil
}
