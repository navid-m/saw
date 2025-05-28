package bingscraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"swarm/models"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeBing(query string) ([]models.SearchResult, error) {
	var (
		results   []models.SearchResult
		resp, err = http.Get("https://www.bing.com/search?q=" + url.QueryEscape(query))
	)

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

	doc.Find("li.b_algo").Each(func(i int, s *goquery.Selection) {
		var (
			titleTag    = s.Find("h2 a")
			descTag     = s.Find(".b_caption p")
			title       = strings.TrimSpace(titleTag.Text())
			link, _     = titleTag.Attr("href")
			description = strings.TrimSpace(descTag.Text())
		)
		if title != "" && link != "" {
			results = append(results, models.SearchResult{
				Title:       title,
				Link:        link,
				Description: description,
			})
		}
	})

	return results, nil
}
