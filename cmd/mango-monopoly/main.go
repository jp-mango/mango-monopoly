package main

import (
	"fmt"
	"log/slog"
	"mango-monopoly/internal/db"
	"mango-monopoly/internal/scraper"
	"mango-monopoly/internal/utils"
	"time"
)

func main() {
	scraper.PullGwinnettAuctionData()
	scraper.DownloadGwinnettAuctionData()
	scraper.PullPauldingAuctionData()
	scraper.DownloadPauldingAuctionData()

	db.ResetDB()

	mango_monopoly, err := db.DbConnect()
	if err != nil {
		slog.Error("failed to connect to the database", "error", err)
	}
	defer mango_monopoly.Close()

	gwinnettPastSales := utils.ReadCSV("Gwinnett-Past-Sales.csv", "Gwinnett")
	pastSalesRowsChanged, err := db.InsertGwinnettPastSalesData(gwinnettPastSales, mango_monopoly)
	if err != nil {
		slog.Error("failed to insert data into db", "error", err)

	}
	fmt.Printf("Past sales rows changed: %d\n\n", pastSalesRowsChanged)

	gwinnettUpcomingSales := utils.ReadCSV("Gwinnett-Upcoming-Sales.csv", "Gwinnett")
	upcomingSalesRowsChanged, err := db.InsertGwinnettUpcomingSalesData(gwinnettUpcomingSales, mango_monopoly)
	if err != nil {
		slog.Error("failed to insert data into db", "error", err)
	}
	fmt.Printf("Upcoming sales rows changed: %d\n\n", upcomingSalesRowsChanged)

	fmt.Printf("New rows in Properties table: %d\n\n", pastSalesRowsChanged+upcomingSalesRowsChanged)

	startTime := time.Now()
	propertiesRowsChanged, err := db.UpdatePropertiesTable_Tax(mango_monopoly)
	if err != nil {
		slog.Error("failed to insert data into db", "error", err)
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Rows updated with tax assessor info in Properties table: %d\n\nProcess took: %v\n\n", propertiesRowsChanged, elapsedTime)
}
