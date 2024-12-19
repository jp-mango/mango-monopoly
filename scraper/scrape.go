package scraper

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
}

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

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			if e.Text == "List of Properties" {
				foundLink = true
				rawLink := e.Attr("href")
				link = fmt.Sprintf("https://%s", strings.TrimPrefix(county.Domain+rawLink, "www."))
				//fmt.Println("Found link:", link)

				scriptPath, err := filepath.Abs("./scraper/main.py")
				if err != nil {
					Logger.Error("error getting absolute path", "err", err)
					return
				}
				//fmt.Println("Script path:", scriptPath)

				pythonCMD := exec.Command("uv", "run", scriptPath, link, county.Name)

				// Ensure the environment PATH is passed
				pythonCMD.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))

				// Capture the combined output of the command
				_, err = pythonCMD.CombinedOutput()
				if err != nil {
					Logger.Error("error running python script", "err", err)
					return
				}

				//fmt.Println("Python script output:", string(output))
			}
		})

		// Debugging for requests
		c.OnRequest(func(r *colly.Request) {
			Logger.Info(fmt.Sprint("Visiting: ", r.URL.String()))
		})

		// Start scraping
		err := c.Visit(county.Webpage)
		if err != nil {
			Logger.Error("error visiting webpage:", "err", err)
			return err
		}

		// Check if the link was found
		if !foundLink {
			Logger.Error("unable to find upcoming property list", "err", err)
			return err
		}

		return nil
	}

	Logger.Error(fmt.Sprintf("unable to find county: '%s'", county.Name))
	return fmt.Errorf("unable to find county: %s", county.Name)
}

func ProcessCSV(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		Logger.Error("unable to read directory", "err", err)
		return nil, err
	}

	if len(files) == 0 {
		Logger.Error("No files found in directory", "err", fmt.Sprintf("directory: %s", files))
		return nil, fmt.Errorf("no files found in directory: %s", path)
	}

	file := files[0]

	fileInfo, err := file.Info()
	if err != nil {
		Logger.Error("unable to retrieve file info", "err", err)
		return nil, err
	}

	filePath := filepath.Join(path, fileInfo.Name())

	f, err := os.Open(filePath)
	if err != nil {
		Logger.Error("unable to open file", "err", err)
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		Logger.Error("unable to read csv", "err", err)
		return nil, err
	}

	var parcelIDs []string
	for index, r := range records {
		if index == 0 {
			continue
		}
		parcelIDs = append(parcelIDs, r[0])
	}

	return parcelIDs, nil
}
