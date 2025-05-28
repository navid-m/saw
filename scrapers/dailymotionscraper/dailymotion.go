package dailymotionscraper

import (
	"fmt"
	"net/http"
	"strings"

	"swarm/models"
	"swarm/models/origins"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
)

func ScrapeDailymotion(query string) ([]models.SearchResult, error) {
	var results []models.SearchResult

	searchURL := "https://www.dailymotion.com/search/" + strings.ReplaceAll(query, " ", "%20") + "/top-results"

	client := &http.Client{}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", uarand.GetRandom())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(`div[data-testid="video-card"]`).Each(func(i int, s *goquery.Selection) {
		var (
			titleTag    = s.Find("a.VideoCard__videoTitleLink___1IRDu")
			link, _     = titleTag.Attr("href")
			title       = strings.TrimSpace(titleTag.Text())
			description = ""
		)

		if !strings.HasPrefix(link, "http") {
			link = "https://www.dailymotion.com" + link
		}

		if title != "" && link != "" {
			results = append(results, models.SearchResult{
				Origin:      origins.DailyMotion,
				Title:       title,
				Link:        link,
				Description: description,
			})
		}
	})

	return results, nil
}
