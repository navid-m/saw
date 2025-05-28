package dailymotionscraper

type dmApiVideo struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

type dmApiResponse struct {
	List []dmApiVideo `json:"list"`
}
