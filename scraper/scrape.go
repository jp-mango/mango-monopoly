package scraper

import (
	"fmt"
	"log/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	if county.Name == "Gwinnett" {
		c := colly.NewCollector(
			colly.AllowedDomains(county.Domain),
		)

		var foundLink bool
		var link string

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			if e.Text == "List of Properties" {
				foundLink = true
				rawLink := e.Attr("href")
				link = fmt.Sprintf("https://%s", strings.TrimPrefix(county.Domain+rawLink, "www."))
				//fmt.Println("Found link:", link)

				scriptPath, err := filepath.Abs("./scraper/main.py")
				if err != nil {
					logger.Error("error getting absolute path", "err", err)
					return
				}
				//fmt.Println("Script path:", scriptPath)

				pythonCMD := exec.Command("uv", "run", scriptPath, link, county.Name)

				// Ensure the environment PATH is passed
				pythonCMD.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))

				// Capture the combined output of the command
				_, err = pythonCMD.CombinedOutput()
				if err != nil {
					logger.Error("error running python script", "err", err)
					return
				}

				//fmt.Println("Python script output:", string(output))
			}
		})

		// Debugging for requests
		c.OnRequest(func(r *colly.Request) {
			logger.Info(fmt.Sprint("Visiting: ", r.URL.String()))
		})

		// Start scraping
		err := c.Visit(county.Webpage)
		if err != nil {
			logger.Error("error visiting webpage:", "err", err)
			return err
		}

		// Check if the link was found
		if !foundLink {
			logger.Error("unable to find upcoming property list", "err", err)
			return err
		}

		return nil
	}

	logger.Error(fmt.Sprintf("unable to find county: '%s'", county.Name))
	return fmt.Errorf("unable to find county: %s", county.Name)
}
