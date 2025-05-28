package googlescraper

// The structure of Google Custom Search API response
type GoogleSearchResponse struct {
	Items []GoogleSearchItem `json:"items"`
}

// A single search result from Google Custom Search API
type GoogleSearchItem struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}
