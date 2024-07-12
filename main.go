package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	downloadGwinnettAuctionData()
}

func downloadGwinnettAuctionData() {
	c := colly.NewCollector(
		colly.AllowedDomains("gwinnetttaxcommissioner.publicaccessnow.com"),
	)

	var upcomingSalesURL string
	var pastResultsURL string
	var upcomingAuctions []string

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
				upcomingAuctions = append(upcomingAuctions, trimmedDate)
			}
		}
	})

	// On request log the URL
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting: %s\n\n", r.URL.String())
	})

	// Visit the page to scrape the link
	gwinnettTaxURL := "https://gwinnetttaxcommissioner.publicaccessnow.com/PropertyTax/DelinquentTax/TaxLiensTaxSales.aspx"
	err := c.Visit(gwinnettTaxURL)
	if err != nil {
		slog.Error(fmt.Sprintf("Error visiting the page: %v", err))
	}

	//Print upcoming auctions
	if len(upcomingAuctions) > 0 {
		fmt.Printf("Upcoming Auction Dates\n")
		fmt.Println("------------------")
		for _, date := range upcomingAuctions {
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

	//Create directory to hold pdfs
	gwinnettDir := "tax-auction/Gwinnett"
	if err := os.MkdirAll(gwinnettDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	// If a PDF link was found, download the PDF
	if pastResultsURL != gwinnettTaxURL && upcomingSalesURL != gwinnettTaxURL {
		// remove old files
		os.Remove("tax-auction/Gwinnett/Gwinnett-Past-Sales.pdf")
		os.Remove("tax-auction/Gwinnett/Gwinnett-Upcoming-Sales.pdf")

		// Create the file paths
		upcomingSalesPDF := "Gwinnett-Upcoming-Sales.pdf"
		upcomingSalesFilepath := filepath.Join(gwinnettDir, upcomingSalesPDF)

		pastResultsPDF := "Gwinnett-Past-Sales.pdf"
		pastResultsFilepath := filepath.Join(gwinnettDir, pastResultsPDF)

		// Download the files
		err := downloadFile(upcomingSalesFilepath, upcomingSalesURL)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		err = downloadFile(pastResultsFilepath, pastResultsURL)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		fmt.Printf("Gwinnett auction data downloaded successfully to: ./%s\n----------------------------------------------------------------------------------------------------------------------------------\n", gwinnettDir)
	} else if pastResultsURL == gwinnettTaxURL {
		os.Remove("tax-auction/Gwinnett/Gwinnett-Upcoming-Sales.pdf")

		upcomingSalesPDF := "Gwinnett-Upcoming-Sales.pdf"
		upcomingSalesFilepath := filepath.Join(gwinnettDir, upcomingSalesPDF)

		err := downloadFile(upcomingSalesFilepath, upcomingSalesURL)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		fmt.Print("No PDF found for past sale history. Upcoming sales updated.\n----------------------------------------------------------------------------------------------------------------------------------\n")
	} else if upcomingSalesURL == gwinnettTaxURL {
		os.Remove("tax-auction/Gwinnett/Gwinnett-Past-Sales.pdf")

		pastResultsPDF := "Gwinnett-Past-Sales.pdf"
		pastResultsFilepath := filepath.Join(gwinnettDir, pastResultsPDF)

		err = downloadFile(pastResultsFilepath, pastResultsURL)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		fmt.Print("No PDF found for upcoming sales. Past sales updated.\n----------------------------------------------------------------------------------------------------------------------------------\n")
	}
}

// Function to download a file from a URL
func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
