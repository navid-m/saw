package dailymotionscraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"swarm/models"
	"swarm/models/origins"
)

func ScrapeDailymotion(query string) ([]models.SearchResult, error) {
	var (
		results []models.SearchResult
		baseURL = "https://api.dailymotion.com/videos"
		params  = url.Values{}
	)

	params.Set("search", query)
	params.Set("limit", "10")
	params.Set("fields", "id,title,duration")

	var (
		apiURL    = baseURL + "?" + params.Encode()
		resp, err = http.Get(apiURL)
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var apiResp dmApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	for _, v := range apiResp.List {
		var (
			link = fmt.Sprintf("https://www.dailymotion.com/video/%s", v.ID)
			desc = fmt.Sprintf("Duration: %d seconds", v.Duration)
		)
		results = append(results, models.SearchResult{
			Origin:      origins.DailyMotion,
			Title:       strings.TrimSpace(v.Title),
			Link:        link,
			Description: desc,
		})
	}

	return results, nil
}
