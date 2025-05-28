package models

// The search engine that yielded some result
type ResultOrigin int

const (
	Bing ResultOrigin = iota
	Google
	DuckDuckGo
	Yahoo
	Qwant
	YouTube
	BitChute
	DailyMotion
	GitHub
	GitLab
	Codeberg
)

// Represents a single search result
type SearchResult struct {
	Origin      ResultOrigin
	Title       string
	Link        string
	Description string
}
