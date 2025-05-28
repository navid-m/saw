package yahooscraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"swarm/models"
	"swarm/models/origins"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
)

func ScrapeYahoo(query string) ([]models.SearchResult, error) {
	var results []models.SearchResult

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://search.yahoo.com/search?p="+url.QueryEscape(query), nil)
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

	doc.Find("li div.dd.algo").Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("div.compTitle > a")
		titleTag := aTag.Find("h3.title")
		descTag := s.Find("div.compText > p")

		title := strings.TrimSpace(titleTag.Text())
		link, _ := aTag.Attr("href")
		description := strings.TrimSpace(descTag.Text())

		if title != "" && link != "" {
			results = append(results, models.SearchResult{
				Origin:      origins.Yahoo,
				Title:       title,
				Link:        link,
				Description: description,
			})
		}
	})

	return results, nil
}
