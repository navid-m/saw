package qwantscraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"swarm/models"

	"github.com/corpix/uarand"
)

type qwantResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result struct {
			Items struct {
				Mainline []struct {
					Type  string `json:"type"`
					Items []struct {
						Title string `json:"title"`
						URL   string `json:"url"`
						Desc  string `json:"desc"`
					} `json:"items"`
				} `json:"mainline"`
			} `json:"items"`
		} `json:"result"`
		ErrorCode int `json:"error_code,omitempty"`
		Message   []string
	} `json:"data"`
}

func ScrapeQwant(query string) ([]models.SearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	const (
		baseURL     = "https://api.qwant.com/v3/search/web"
		resultCount = 10
		locale      = "en_US"
		safesearch  = "1"
	)

	params := url.Values{}
	params.Set("q", query)
	params.Set("count", fmt.Sprintf("%d", resultCount))
	params.Set("offset", "0")
	params.Set("locale", locale)
	params.Set("safesearch", safesearch)

	apiURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", uarand.GetRandom())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	var apiResp qwantResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status != "success" {
		return nil, fmt.Errorf("api error: code %d, message %v", apiResp.Data.ErrorCode, apiResp.Data.Message)
	}

	results := make([]models.SearchResult, 0)

	for _, mainlineBlock := range apiResp.Data.Result.Items.Mainline {
		if mainlineBlock.Type != "web" {
			continue
		}

		for _, item := range mainlineBlock.Items {
			if item.Title != "" && item.URL != "" {
				results = append(results, models.SearchResult{
					Title:       item.Title,
					Link:        item.URL,
					Description: item.Desc,
				})
			}
		}
	}

	return results, nil
}
