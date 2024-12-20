package scraper

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func (county *CountyScraper) ScrapeAuctionData() error {
	switch county.Name {
	case "Gwinnett":
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
			Logger.Info(fmt.Sprintf("Visiting: %s", r.URL.String()))
		})

		c.OnResponse(func(r *colly.Response) {
			Logger.Info(fmt.Sprintf("Status: %d", r.StatusCode))
		})

		c.OnError(func(r *colly.Response, err error) {
			Logger.Error(fmt.Sprintf("request URL: %s failed with response: %v", r.Request.URL, r), "err", err)
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
	default:
		Logger.Error(fmt.Sprintf("unable to find county: '%s'", county.Name))
		return fmt.Errorf("unable to find county: %s", county.Name)
	}
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

type ParcelData struct {
	PropertyID    string
	AlternateID   string
	Address       string
	PropertyClass string
	Neighborhood  int32
	DeedAcres     float32
}

func ScrapeParcelData(parcelIDs []string) error {
	c := colly.NewCollector(
		colly.AllowedDomains("gwinnettassessor.manatron.com"),
	)

	for i := 0; i < len(parcelIDs); i++ {
		parcelIDs[i] = strings.Replace(parcelIDs[i], " ", "%20", -1)

		url := fmt.Sprintf("https://gwinnettassessor.manatron.com/IWantTo/PropertyGISSearch/PropertyDetail.aspx?p=%s", parcelIDs[i])

		// Set up a slice to hold the scraped data
		var parcelData ParcelData

		// Scrape the content of the relevant `div`
		c.OnHTML("div#dnn_ctr1385_ContentPane", func(e *colly.HTMLElement) {
			// Extract table rows within the div
			e.ForEach("table tr", func(_ int, row *colly.HTMLElement) {
				header := strings.TrimSpace(row.ChildText("th"))
				value := strings.TrimSpace(row.ChildText("td"))

				switch header {
				case "Property ID":
					parcelData.PropertyID = strings.TrimSpace(value)
				case "Alternate ID":
					parcelData.AlternateID = strings.TrimSpace(value)
				case "Address":
					parcelData.Address = strings.TrimSpace(value)
				case "Property Class":
					parcelData.PropertyClass = strings.TrimSpace(value)
				case "Neighborhood":
					n, err := strconv.Atoi(strings.TrimSpace(value))
					if err != nil {
						Logger.Error("unable to convert to int", "err", err)
					}
					parcelData.Neighborhood = int32(n)
				case "Deed Acres":
					da, err := strconv.ParseFloat(strings.TrimSpace(value), 32)
					if err != nil {
						Logger.Error("unable to convert to float", "err", err)
					}
					parcelData.DeedAcres = float32(da)
				}
			})
		})
		//TODO: scrape remainign data from page and load into DB

		c.OnResponse(func(r *colly.Response) {
			if r.StatusCode != 200 {
				fmt.Printf("Status: %d\n", r.StatusCode)
			}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Printf("Request URL: %s failed with response: %v\nError: %v\n", r.Request.URL, r, err)
		})

		err := c.Visit(url)
		if err != nil {
			fmt.Printf("Error visiting webpage: %v\n", err)
			return err
		}

		// Print the scraped data
		fmt.Printf("\nScraped Parcel Data: %+v\n", parcelData)
	}

	return nil
}
