package qwantscraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"swarm/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
)

func ScrapeQwant(query string) ([]models.SearchResult, error) {
	var (
		client   = &http.Client{}
		reqUrl   = "https://www.qwant.com/?q=" + url.QueryEscape(query) + "&t=web"
		req, err = http.NewRequest("GET", reqUrl, nil)
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

	results, failed, result, err := extractFromQwant(resp, nil)
	if failed {
		return result, err
	}

	return results, nil
}

func extractFromQwant(resp *http.Response, results []models.SearchResult) (
	[]models.SearchResult,
	bool,
	[]models.SearchResult,
	error,
) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, true, nil, err
	}

	doc.Find(`div[data-testid="SERVariant-A"]`).Each(func(i int, s *goquery.Selection) {
		var (
			linkTag     = s.Find(`a.external`).First()
			titleTag    = s.Find(`div.HhS7p span`).First()
			descTag     = s.Find(`div.IXeY3.aVNer.Zxaxi.nSmMx`).First()
			title       = strings.TrimSpace(titleTag.Text())
			link, _     = linkTag.Attr("href")
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

	return results, false, nil, nil
}
