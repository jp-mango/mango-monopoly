package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
)

func main() {
	downloadGwinnettAuctionData()
	readPDF()
}

func readPDF() {
	gwinnettSaleHistoryPath := "tax-auction/Gwinnett/Gwinnett-Past-Sales.pdf"
	gsh, err := os.Open(gwinnettSaleHistoryPath)
	if err != nil {
		slog.Error(fmt.Sprintf("unable to open file: %v", err))
	}
	defer gsh.Close()
}

func downloadGwinnettAuctionData() {
	c := colly.NewCollector(
		colly.AllowedDomains("gwinnetttaxcommissioner.publicaccessnow.com"),
	)

	var upcomingSales string
	var pastResults string

	// Select div - #: ID, .:Class
	// Gwinnett's upcoming sales
	c.OnHTML("#dnn_ctr1334_ModuleContent a", func(us *colly.HTMLElement) {
		upcomingSales = us.Request.AbsoluteURL(us.Attr("href"))
		fmt.Printf("Gwinnett's upcoming auctions: %s\n\n", upcomingSales)
	})

	// Gwinnett's past results
	c.OnHTML("#dnn_ctr1341_ContentPane a", func(pr *colly.HTMLElement) {
		pastResults = pr.Request.AbsoluteURL(pr.Attr("href"))
		fmt.Printf("Gwinnett's auction sale history: %s\n\n", pastResults)
	})

	// On request log the URL
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting\n\n", r.URL.String())
	})

	// Visit the page to scrape the link
	err := c.Visit("https://gwinnetttaxcommissioner.publicaccessnow.com/PropertyTax/DelinquentTax/TaxLiensTaxSales.aspx")
	if err != nil {
		log.Fatalf("Error visiting the page: %v", err)
	}

	// If a PDF link was found, download the PDF
	if pastResults != "" && upcomingSales != "" {
		// empty directory content
		os.RemoveAll("tax-auction/Gwinnett")
		// Ensure the directory exists
		gwinnettDir := "tax-auction/Gwinnett"
		if err := os.MkdirAll(gwinnettDir, os.ModePerm); err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}

		// Create the file paths
		upcomingSalesPDF := "Gwinnett-Upcoming-Sales.pdf"
		upcomingSalesFilepath := filepath.Join(gwinnettDir, upcomingSalesPDF)

		pastResultsPDF := "Gwinnett-Past-Sales.pdf"
		pastResultsFilepath := filepath.Join(gwinnettDir, pastResultsPDF)

		// Download the files
		err := downloadFile(upcomingSalesFilepath, upcomingSales)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}

		err = downloadFile(pastResultsFilepath, pastResults)
		if err != nil {
			log.Fatalf("Error downloading upcoming sales: %v", err)
		}
		fmt.Print("Gwinnett auction data downloaded successfully!\n----------------------------------------------------------------------------------------------------------------------------------\n")
	} else {
		fmt.Println("No PDF link found.")
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
