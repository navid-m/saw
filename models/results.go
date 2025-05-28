package models

import "github.com/navid-m/saw/models/origins"

// Represents a single search result
type SearchResult struct {
	Origin      origins.ResultOrigin
	Title       string
	Link        string
	Description string
}
