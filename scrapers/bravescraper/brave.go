package bravescraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"swarm/models"
	"swarm/models/origins"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeBrave(query string) ([]models.SearchResult, error) {
	var (
		results   []models.SearchResult
		resp, err = http.Get("https://search.brave.com/search?q=" + url.QueryEscape(query))
	)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	results, failed, result, err := extractFromBrave(resp, results)
	if failed {
		return result, err
	}

	return results, nil
}

func extractFromBrave(resp *http.Response, results []models.SearchResult) (
	[]models.SearchResult,
	bool,
	[]models.SearchResult,
	error,
) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, true, nil, err
	}

	doc.Find("div[data-type='web']").Each(func(i int, s *goquery.Selection) {
		var (
			titleTag    = s.Find(".title")
			linkTag     = s.Find("a[href]").First()
			descTag     = s.Find(".snippet-description")
			title       = strings.TrimSpace(titleTag.Text())
			link, _     = linkTag.Attr("href")
			description = strings.TrimSpace(descTag.Text())
		)

		if title != "" && link != "" {
			results = append(results, models.SearchResult{
				Origin:      origins.Brave,
				Title:       title,
				Link:        link,
				Description: description,
			})
		}
	})

	// Alternative selector pattern for Brave search results
	if len(results) == 0 {
		doc.Find("div.snippet").Each(func(i int, s *goquery.Selection) {
			var (
				titleTag    = s.Find(".title")
				linkTag     = s.Find("a").First()
				descTag     = s.Find(".snippet-description")
				title       = strings.TrimSpace(titleTag.Text())
				link, _     = linkTag.Attr("href")
				description = strings.TrimSpace(descTag.Text())
			)

			if title != "" && link != "" {
				results = append(results, models.SearchResult{
					Origin:      origins.Brave,
					Title:       title,
					Link:        link,
					Description: description,
				})
			}
		})
	}

	return results, false, nil, nil
}
