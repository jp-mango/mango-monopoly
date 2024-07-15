package main

import (
	"fmt"
	"log"
	"mango-monopoly/internal/db"
	"mango-monopoly/internal/utils"
)

func main() {
	//scraper.DownloadGwinnettAuctionData()

	//db.ResetDB()

	mango_monopoly, err := db.DbConnect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer mango_monopoly.Close()

	gwinnettPastSales := utils.ReadCSV("Gwinnett-Past-Sales.csv", "Gwinnett")
	result, err := db.InsertGwinnettPastSalesData(gwinnettPastSales, mango_monopoly)
	if err != nil {
		log.Fatalf("failed to insert data into db: %v", err)
	}
	fmt.Printf("Rows Affected: %d\n", result)

	gwinnettUpcomingSales := utils.ReadCSV("Gwinnett-Upcoming-Sales.csv", "Gwinnett")
	result, err = db.InsertGwinnettUpcomingSalesData(gwinnettUpcomingSales, mango_monopoly)
	if err != nil {
		log.Fatalf("failed to insert data into db: %v", err)
	}
	fmt.Printf("Rows Affected: %d", result)

}
