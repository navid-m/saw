package main

import (
	"fmt"
	"log"
	"swarm/scrapers/yahooscraper"
)

func main() {
	query := "asdf"
	// results, err := bingscraper.ScrapeBing(query)
	results, err := yahooscraper.ScrapeYahoo(query)
	if err != nil {
		log.Fatalf("Error scraping Bing: %v", err)
	}
	fmt.Println(len(results))
	for i, r := range results {
		fmt.Printf("Result #%d\n", i+1)
		fmt.Printf("Title: %s\n", r.Title)
		fmt.Printf("Link: %s\n", r.Link)
		fmt.Printf("Description: %s\n\n", r.Description)
	}
}
