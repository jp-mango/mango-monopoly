package scraper

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log/slog"
	"mango-monopoly/internal/models"
	"math"
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
				output, err := pythonCMD.CombinedOutput()
				if err != nil {
					Logger.Error("Python script failed", "output", string(output), "err", err)
				} else {
					Logger.Info("Python script succeeded", "output", string(output))
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

func ScrapeGwinnettParcelData(parcelIDs []string) error {
	c := colly.NewCollector(
		colly.AllowedDomains("gwinnettassessor.manatron.com"),
	)

	for i := 0; i < len(parcelIDs); i++ {
		var prop models.Property

		parcelIDs[i] = strings.Replace(parcelIDs[i], " ", "%20", -1)

		url := fmt.Sprintf("https://gwinnettassessor.manatron.com/IWantTo/PropertyGISSearch/PropertyDetail.aspx?p=%s", parcelIDs[i])
		prop.TaxURL = sql.NullString{String: url, Valid: true}

		state := "Georgia"
		prop.State = sql.NullString{String: state, Valid: true}

		// Scrape the content of the relevant `div`
		c.OnHTML("div#dnn_ctr1385_ContentPane", func(e *colly.HTMLElement) {
			// Extract table rows within the div
			e.ForEach("table tr", func(_ int, row *colly.HTMLElement) {
				header := strings.TrimSpace(row.ChildText("th"))
				value := strings.TrimSpace(row.ChildText("td"))

				switch header {
				case "Property ID":
					prop.ParcelID = sql.NullString{String: value, Valid: true}
				case "Address":
					prop.Address = sql.NullString{String: value, Valid: true}
				case "Property Class":
					prop.PropertyClass = sql.NullString{String: value, Valid: true}
				case "Deed Acres":
					da, err := strconv.ParseFloat(strings.TrimSpace(value), 32)
					if err != nil {
						Logger.Error("unable to convert to float", "err", err)
					}
					prop.LotSize = sql.NullFloat64{Float64: math.Round(da*100) / 100, Valid: true}
				}
			})
		})

		c.OnHTML("div#lxT1388", func(e *colly.HTMLElement) {
			e.ForEach("table.ui-widget-content.ui-table tr", func(_ int, row *colly.HTMLElement) {
				header := strings.TrimSpace(row.ChildText("th"))
				value := strings.TrimSpace(row.ChildText("td"))

				switch header {
				case "Type":
					prop.PropertyType = sql.NullString{String: value, Valid: true}
				case "Grade":
					prop.Grade = sql.NullString{String: value, Valid: true}
				case "Year Built":
					yb, err := strconv.ParseInt(value, 10, 16)
					if err != nil {
						Logger.Error("unable to convert to int", "err", err)
					}
					prop.YearBuilt = sql.NullInt16{Int16: int16(yb), Valid: true}
				}
			})
		})

		c.OnHTML("table#ValueHistory", func(e *colly.HTMLElement) {

			e.ForEach("tr", func(rowIndex int, row *colly.HTMLElement) {
				// Handle the first row (header)
				if rowIndex == 0 {
					return
				}

				// Extract data rows (rowIndex > 0)
				attribute := strings.TrimSpace(row.ChildText("th"))
				if attribute == "" {
					return
				}

				// Only extract values for the most current year (second column)
				value := strings.TrimSpace(row.ChildText("td:nth-of-type(1)")) // First column after the header
				switch attribute {
				case "Land Val":
					lv, err := strconv.ParseInt(strings.ReplaceAll(strings.ReplaceAll(value, "$", ""), ",", ""), 10, 64)
					if err != nil {
						Logger.Error("unable to parse land value", "value", value, "err", err)
					}
					prop.LandValue = sql.NullInt64{Int64: lv, Valid: true}
				case "Imp Val":
					iv, err := strconv.ParseInt(strings.ReplaceAll(strings.ReplaceAll(value, "$", ""), ",", ""), 10, 64)
					if err != nil {
						Logger.Error("unable to parse improvement value", "value", value, "err", err)
					}
					prop.ImprovementValue = sql.NullInt64{Int64: iv, Valid: true}
				case "Total Appr":
					ta, err := strconv.ParseInt(strings.ReplaceAll(strings.ReplaceAll(value, "$", ""), ",", ""), 10, 64)
					if err != nil {
						Logger.Error("unable to parse total appraisal", "value", value, "err", err)
					}
					prop.AppraisalValue = sql.NullInt64{Int64: ta, Valid: true}
				}
			})
		})

		c.OnHTML("div#1388Attributes", func(e *colly.HTMLElement) {
			// Focus on the table body rows
			e.ForEach("table#Attribute tbody tr", func(rowIndex int, row *colly.HTMLElement) {
				// Skip the first "jqgfirstrow" which is often just a spacer
				if row.Attr("class") == "jqgfirstrow" {
					return
				}

				// The attribute name is in the second td
				attribute := strings.TrimSpace(row.ChildText("td:nth-of-type(2)"))
				// The detail is in the third td
				detail := strings.TrimSpace(row.ChildText("td:nth-of-type(3)"))

				switch attribute {
				case "Roof Structure":
					prop.RoofStructure = sql.NullString{String: detail, Valid: true}
				case "Roof Cover":
					prop.RoofCover = sql.NullString{String: detail, Valid: true}
				case "Heating":
					prop.Heating = sql.NullString{String: detail, Valid: true}
				case "A/C":
					prop.Cooling = sql.NullString{String: detail, Valid: true}
				case "Stories":
					s, err := strconv.ParseFloat(detail, 64)
					if err != nil {
						Logger.Error("unable to convert stories to float", "err", err)
					}
					prop.Floors = sql.NullFloat64{Float64: s, Valid: true}
				case "Bedrooms":
					b, err := strconv.Atoi(detail)
					if err != nil {
						Logger.Error("unable to convert bedrooms to int", "err", err)
					}
					prop.Bedrooms = sql.NullInt16{Int16: int16(b), Valid: true}
				case "Bathrooms":
					b, err := strconv.ParseFloat(detail, 64)
					if err != nil {
						Logger.Error("unable to convert bathrooms to float", "err", err)
					}
					prop.Bathrooms = sql.NullFloat64{Float64: b, Valid: true}
				}
			})
		})

		c.OnHTML("img#sketch", func(e *colly.HTMLElement) {
			imgURL := e.Attr("src")
			if imgURL != "" {
				// Convert to absolute URL if necessary
				imgURL = e.Request.AbsoluteURL(imgURL)
				prop.FloorPlanPhoto = sql.NullString{String: imgURL, Valid: true}
			}
		})

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
		fmt.Printf("\nScraped Parcel Data: %+v\n", prop)
	}

	return nil
}
