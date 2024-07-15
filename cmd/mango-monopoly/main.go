package main

import (
	"log"
	"mango-monopoly/internal/db"
	"mango-monopoly/internal/utils"
)

func main() {
	//scraper.DownloadGwinnettAuctionData()

	//db.DbConnect()
	db.ResetDB()

	mango_monopoly, err := db.DbConnect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer mango_monopoly.Close()

	gwinnettPastSales := utils.ReadCSV("Gwinnett-Past-Sales.csv")

	err = db.InsertGwinnettPastSalesData(gwinnettPastSales, mango_monopoly)
	if err != nil {
		log.Fatalf("failed to insert data into db: %v", err)
	}
}
