package scraper

import (
	"fmt"
	"log"
	"log/slog"
	"mango-monopoly/internal/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

var (
	PauldingTaxURL        = "http://www.paulding.gov/208/Tax-Commissioner"
	pauldingUpcomingSales string
)

func PullPauldingAuctionData() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.paulding.gov", "paulding.gov"),
	)

	c.OnHTML("#SideItem1149 a", func(al *colly.HTMLElement) {
		pauldingUpcomingSalesPage := al.Request.AbsoluteURL(al.Attr("href"))
		if pauldingUpcomingSalesPage != "" {
			// Use a new collector instance to avoid URL already visited error
			c2 := c.Clone()
			c2.OnHTML("#div41b4b299-0ca0-4644-a522-3eb1547321bf a", func(e *colly.HTMLElement) {
				documentURL := e.Attr("href")
				if strings.Contains(documentURL, "/DocumentCenter/View/") {
					pauldingUpcomingSales = e.Request.AbsoluteURL(documentURL)
					fmt.Printf("Found document URL: %s\n", pauldingUpcomingSales)
				}
			})

			err := c2.Visit(pauldingUpcomingSalesPage)
			if err != nil {
				slog.Error(fmt.Sprintf("Error visiting navbar link: %v", err))
			}
		} else {
			slog.Info("No navbar item found with ID: SideItem1149")
		}
	})

	// On request log the URL
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting: %s\n", r.URL.String())
	})

	// Visit the initial page
	err := c.Visit(PauldingTaxURL)
	if err != nil {
		slog.Error(fmt.Sprintf("Error visiting the page: %v", err))
	}
}

func DownloadPauldingAuctionData() {
	pauldingDir := filepath.Join("tax-auction", "Paulding")
	if err := os.MkdirAll(pauldingDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	if pauldingUpcomingSales == "" { // figure out if link is incorrect
		log.Fatal("unable to find link to Paulding's upcoming sales")
	}

	upcomingSalesPDF := filepath.Join(pauldingDir, "pdf", "Paulding-Upcoming-Sales.pdf")

	os.Remove(upcomingSalesPDF)

	err := utils.DownloadFile(upcomingSalesPDF, pauldingUpcomingSales)
	if err != nil {
		log.Fatalf("Error downloading Paulding upcoming sales: %v", err)
	}

	utils.RunPdfExtraction(filepath.Join("scripts", "paulding-pdf-extract.py"), "Paulding")

	fmt.Printf("Paulding auction data downloaded successfully to: .\\%s\n------------------------------------------------------------------------------\n", pauldingDir)
}

type PauldingPropertyDetails struct {
	Account string
	Owner   string
	Address string
	Zip     string
	Updates int64
}

func PauldingTaxRecordsPull(realKey string) {
	c := colly.NewCollector()

	//propertyURL := fmt.Sprintf("https://qpublic.schneidercorp.com/Application.aspx?App=PauldingCountyGA&Layer=Parcels&PageType=Search")

	details := PauldingPropertyDetails{}

	c.OnResponse(func(r *colly.Response) {
		details.Updates++
	})
}
