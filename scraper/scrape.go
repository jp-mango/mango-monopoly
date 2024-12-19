package scraper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

type Scraper interface {
	Scrape() error
}

type CountyScraper struct {
	Name    string
	Webpage string
	Domain  string
}

func (county *CountyScraper) Scrape() error {
	if county.Name == "Gwinnett" {
		c := colly.NewCollector(
			colly.AllowedDomains(county.Domain),
		)

		var foundLink bool
		var link string
		var pythonError error

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			if e.Text == "List of Properties" {
				foundLink = true
				rawLink := e.Attr("href")
				link = fmt.Sprintf("https://%s", strings.TrimPrefix(county.Domain+rawLink, "www."))
				//fmt.Println("Found link:", link)

				scriptPath, err := filepath.Abs("./scraper/main.py")
				if err != nil {
					pythonError = fmt.Errorf("error getting absolute path: %w", err)
					return
				}
				//fmt.Println("Script path:", scriptPath)

				pythonCMD := exec.Command("uv", "run", scriptPath, link, county.Name)

				// Ensure the environment PATH is passed
				pythonCMD.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))

				// Capture the combined output of the command
				_, err = pythonCMD.CombinedOutput()
				if err != nil {
					//fmt.Printf("Python script error: %v\nOutput: %s\n", err, string(output))
					pythonError = fmt.Errorf("error running Python script: %w", err)
					return
				}

				//fmt.Println("Python script output:", string(output))
			}
		})

		// Debugging for requests
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		// Start scraping
		err := c.Visit(county.Webpage)
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

	return fmt.Errorf("unable to find county: %s", county.Name)
}
