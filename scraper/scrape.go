package scraper

import (
	"fmt"
	"os/exec"
	"path/filepath"
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

func (gco *GwinnettScraper) Scrape() error {
	c := colly.NewCollector(
		colly.AllowedDomains(gco.Domain),
	)

	var foundLink bool
	var link string
	var pythonError error // To capture Python script errors

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Text == "List of Properties" {
			foundLink = true
			link = strings.SplitN(gco.Domain+e.Attr("href"), "?", 2)[0]
			fmt.Println("Found link:", link)

			//TODO: implement working command below into Go script
			//TODO: python3 ./scraper/main.py https://www.gwinnetttaxcommissioner.com/documents/d/egov/decembertaxsalelist Gwinnett
			scriptPath, err := filepath.Abs("./scraper/main.py")
			if err != nil {
				fmt.Printf("Error getting absolute path: %v\n", err)
				pythonError = fmt.Errorf("error getting absolute path: %w", err)
				return
			}

			// Execute the Python script
			pythonCMD := exec.Command("python3", scriptPath, link, "Gwinnett")
			output, err := pythonCMD.CombinedOutput()
			if err != nil {
				fmt.Printf("Python script error: %v\n", err)
				pythonError = fmt.Errorf("error running Python script: %w", err)
				return
			}
			fmt.Println("Python script output:", string(output))
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

	// Check if the link was found
	if !foundLink {
		return fmt.Errorf("unable to find upcoming property list")
	}

	// Check if there was a Python script error
	if pythonError != nil {
		return pythonError
	}

	return nil
}
