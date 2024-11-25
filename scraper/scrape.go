package scraper

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gocolly/colly"
)

type Scraper interface {
	Scrape() error
}

type GwinnettScraper struct {
	Webpage string
	Domain  string
}

// func (gco *GwinnettScraper) Scrape() error {
func (gco *GwinnettScraper) Scrape() error {
	c := colly.NewCollector(
		colly.AllowedDomains(gco.Domain),
	)

	var foundLink bool

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Text == "List of Properties" {
			foundLink = true

			link := strings.SplitN(gco.Domain+e.Attr("href"), "?", 2)
			fmt.Println("Found link:", link[0])

			pythonCMD := exec.Command("uv", "run", "./scraper/main.py")
			output, err := pythonCMD.CombinedOutput()
			if err != nil {
				fmt.Println("Error running python script:", err)
				return
			}
			fmt.Println(string(output))
		}
	})

	// Debugging for requests
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	err := c.Visit(gco.Webpage)
	if err != nil {
		return fmt.Errorf("error visiting webpage: %w", err)
	}

	// After scraping is done, check if the link was found
	if !foundLink {
		return fmt.Errorf("unable to find upcoming property list")
	}

	return nil
}

func ScrapeCounty(county string) error {
	var scraper Scraper

	switch county {
	case "Gwinnett":
		scraper = &GwinnettScraper{
			Webpage: "",
			Domain:  "",
		}
	default:
		return fmt.Errorf("unknown county: %s", county)
	}

	return scraper.Scrape()
}
