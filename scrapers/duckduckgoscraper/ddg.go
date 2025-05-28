package duckduckgoscraper

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

func ScrapeDuckDuckGo(query string) ([]models.SearchResult, error) {
	var (
		results  []models.SearchResult
		client   = &http.Client{}
		req, err = http.NewRequest("GET", "https://html.duckduckgo.com/html/?q="+url.QueryEscape(query), nil)
	)

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

	doc.Find("div.result.results_links.results_links_deep.web-result").Each(func(i int, s *goquery.Selection) {
		var (
			titleTag    = s.Find("h2.result__title a.result__a")
			descTag     = s.Find("a.result__snippet")
			title       = strings.TrimSpace(titleTag.Text())
			link, _     = titleTag.Attr("href")
			description = strings.TrimSpace(descTag.Text())
		)
		if title != "" && link != "" {
			results = append(results, models.SearchResult{
				Origin:      origins.DuckDuckGo,
				Title:       title,
				Link:        link,
				Description: description,
			})
		}
	})

	return results, nil
}
