package main

import (
	"fmt"
	"log/slog"
	"mango-monopoly/internal/db"
	"mango-monopoly/internal/scraper"
	"mango-monopoly/internal/utils"
)

func main() {
	scraper.PullGwinnettAuctionData()
	scraper.DownloadGwinnettAuctionData()

	//db.ResetDB()

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

	propertiesRowsChanged, err := db.UpdatePropertiesTable(mango_monopoly)
	if err != nil {
		slog.Error("failed to insert data into db", "error", err)

	}
	fmt.Printf("Rows changed in Properties table: %d\n\n", propertiesRowsChanged)
}
