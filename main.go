package main

import (
	"fmt"
	"log"

	"github.com/navid-m/saw/models"
	"github.com/navid-m/saw/scrapers/bingscraper"
	"github.com/navid-m/saw/scrapers/bravescraper"
	"github.com/navid-m/saw/scrapers/dailymotionscraper"
	"github.com/navid-m/saw/scrapers/duckduckgoscraper"
	"github.com/navid-m/saw/scrapers/qwantscraper"
	"github.com/navid-m/saw/scrapers/yahooscraper"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	scrapers := map[string]func(string) ([]models.SearchResult, error){
		"Bing": func(q string) ([]models.SearchResult, error) {
			return bingscraper.ScrapeBing(q)
		},
		"Yahoo": func(q string) ([]models.SearchResult, error) {
			return yahooscraper.ScrapeYahoo(q)
		},
		"Qwant": func(q string) ([]models.SearchResult, error) {
			return qwantscraper.ScrapeQwant(q)
		},
		"DailyMotion": func(q string) ([]models.SearchResult, error) {
			return dailymotionscraper.ScrapeDailymotion(q)
		},
		"DuckDuckGo": func(q string) ([]models.SearchResult, error) {
			return duckduckgoscraper.ScrapeDuckDuckGo(q)
		},
		"Brave": func(q string) ([]models.SearchResult, error) {
			return bravescraper.ScrapeBrave(q)
		},
	}

	var (
		scraperNames    = []string{"Bing", "Yahoo", "Qwant", "DailyMotion", "DuckDuckGo", "Brave"}
		form            = tview.NewForm()
		output          = tview.NewTextView().SetDynamicColors(true).SetWrap(true)
		selectedScraper string
		query           string
	)

	output.SetBorder(true).SetTitle("Results")

	form.AddDropDown("Scraper", scraperNames, 0, func(option string, index int) {
		selectedScraper = option
	})
	form.AddInputField("Query", "", 30, nil, func(text string) {
		query = text
	})
	form.AddButton("Search", func() {
		output.Clear()
		if selectedScraper == "" || query == "" {
			fmt.Fprintln(output, "[red]Please select a scraper and enter a query.")
			return
		}

		results, err := scrapers[selectedScraper](query)
		if err != nil {
			fmt.Fprintf(output, "[red]Error using %s: %v\n", selectedScraper, err)
			return
		}

		if len(results) == 0 {
			fmt.Fprintln(output, "[yellow]No results found.")
			return
		}

		for i, r := range results {
			fmt.Fprintf(output, "[green]Result #%d\n", i+1)
			fmt.Fprintf(output, "Title: [white]%s\n", r.Title)
			fmt.Fprintf(output, "Link: [blue]%s\n", r.Link)
			fmt.Fprintf(output, "Description: [white]%s\n\n", r.Description)
		}
	})

	form.AddButton("Quit", func() {
		app.Stop()
	})

	form.SetBorder(true).SetTitle("SAW").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(form, 0, 1, true).
		AddItem(output, 0, 2, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
