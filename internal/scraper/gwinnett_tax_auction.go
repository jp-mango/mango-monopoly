package scraper

import (
	"fmt"
	"log"
	"log/slog"
	"mango-monopoly/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var gwinnettTaxURL string = "https://gwinnetttaxcommissioner.publicaccessnow.com/PropertyTax/DelinquentTax/TaxLiensTaxSales.aspx"
var GwinnettUpcomingAuctions []string
var upcomingSalesURL string
var pastResultsURL string

func PullGwinnettAuctionData() {
	c := colly.NewCollector(
		colly.AllowedDomains("gwinnetttaxcommissioner.publicaccessnow.com"),
	)

	// Select div - #: ID, .:Class
	// Gwinnett's upcoming sales
	c.OnHTML("#dnn_ctr1334_ModuleContent a", func(us *colly.HTMLElement) {
		upcomingSalesURL = us.Request.AbsoluteURL(us.Attr("href"))
	})

	// Gwinnett's past results
	c.OnHTML("#dnn_ctr1341_ContentPane a", func(pr *colly.HTMLElement) {
		pastResultsURL = pr.Request.AbsoluteURL(pr.Attr("href"))
	})

	//Upcoming auction dates
	c.OnHTML("#dnn_ctr1334_ModuleContent span[style*='font-size: 18px']", func(ad *colly.HTMLElement) {
		dates := strings.Split(ad.Text, "\n")
		for _, date := range dates {
			trimmedDate := strings.TrimSpace(date)
			if trimmedDate != "" {
				GwinnettUpcomingAuctions = append(GwinnettUpcomingAuctions, trimmedDate)
			}
		}
	})

	// On request log the URL
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting: %s\n\n", r.URL.String())
	})

	// Visit the page to scrape the link
	err := c.Visit(gwinnettTaxURL)
	if err != nil {
		slog.Error(fmt.Sprintf("Error visiting the page: %v", err))
	}

	//Print upcoming auctions
	if len(GwinnettUpcomingAuctions) > 0 {
		fmt.Printf("Upcoming Auction Dates\n")
		fmt.Println("------------------")
		for _, date := range GwinnettUpcomingAuctions {
			auctionDate, err := time.Parse("January 2, 2006", date)
			if err != nil {
				slog.Error(fmt.Sprintf("Error parsing date '%s': %v", date, err))
			}

			nextAuctionTime := time.Until(auctionDate)

			days := int(nextAuctionTime.Hours()) / 24
			fmt.Printf("%s | %d days\n", date, days)
		}
		fmt.Println()
	}

}

func DownloadGwinnettAuctionData() {
	//Create directory to hold pdfs
	gwinnettDir := filepath.Join("tax-auction", "Gwinnett")
	if err := os.MkdirAll(gwinnettDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	// If a PDF link was found, download the PDF
	if pastResultsURL != gwinnettTaxURL && upcomingSalesURL != gwinnettTaxURL {
		// Create the file paths
		upcomingSalesPDF := filepath.Join(gwinnettDir, "pdf", "Gwinnett-Upcoming-Sales.pdf")
		pastResultsPDF := filepath.Join(gwinnettDir, "pdf", "Gwinnett-Past-Sales.pdf")

		// remove old files
		os.Remove(upcomingSalesPDF)
		os.Remove(pastResultsPDF)

		// Download the files
		err := utils.DownloadFile(upcomingSalesPDF, upcomingSalesURL)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		err = utils.DownloadFile(pastResultsPDF, pastResultsURL)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		utils.RunPdfExtraction(filepath.Join("scripts", "pdf-extract.py"))

		fmt.Printf("Gwinnett auction data downloaded successfully to: .\\%s\n----------------------------------------------------------------------------------------------------------------------------------\n", gwinnettDir)
	} else if pastResultsURL == gwinnettTaxURL {
		upcomingSalesPDF := "/pdf/Gwinnett-Upcoming-Sales.pdf"
		upcomingSalesFilepath := filepath.Join(gwinnettDir, upcomingSalesPDF)

		os.Remove(upcomingSalesFilepath)

		err := utils.DownloadFile(upcomingSalesFilepath, upcomingSalesURL)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		fmt.Print("No PDF found for past sale history. Upcoming sales pdf updated.\n----------------------------------------------------------------------------------------------------------------------------------\n")
	} else if upcomingSalesURL == gwinnettTaxURL {
		pastResultsPDF := "Gwinnett-Past-Sales.pdf"
		pastResultsFilepath := filepath.Join(gwinnettDir, pastResultsPDF)

		os.Remove(pastResultsFilepath)

		err := utils.DownloadFile(pastResultsFilepath, pastResultsURL)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		fmt.Print("No PDF found for upcoming sales. Past sales pdf updated.\n----------------------------------------------------------------------------------------------------------------------------------\n")
	}
}
