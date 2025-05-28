package origins

// The search engine that yielded some result
type ResultOrigin int

const (
	Bing ResultOrigin = iota
	Google
	DuckDuckGo
	Yahoo
	Qwant
	DailyMotion
)
