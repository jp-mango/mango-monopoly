package main

import (
	"fmt"
	"log"
	"mango-monopoly/internal/db"
	"mango-monopoly/internal/utils"
	"slices"
	"strings"
)

func main() {
	//scraper.DownloadGwinnettAuctionData()

	//db.DbConnect()
	//db.ResetDB()

	database, err := db.DbConnect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	salesData := utils.ReadCSV("Gwinnett-Past-Sales.csv")

	header := salesData[0]
	phrase := "tentatively schduled tax sale"

	for i, value := range salesData {
		if slices.Compare(value, header) != 0 && strings.ToLower(strings.TrimSpace(value[2])) != phrase {
			auctionDate := strings.TrimSpace(value[0])
			parcelID := strings.TrimSpace(value[1])
			previousOwner := strings.TrimSpace(value[2])
			addr := strings.TrimSpace(value[3])
			startingBid := strings.TrimSpace(value[5])
			taxDeedPurchaser := strings.TrimSpace(value[6])
			winningBidAmount := strings.TrimSpace(value[7])

			fmt.Printf("index: %d\n, date: %s\n, PID: %s\n, Previous Owner: %s\n, Address: %s\n, Starting Bid: %s\n, New Owner: %s\n, Amount Paid: %s\n\n", i, auctionDate, parcelID, previousOwner, addr, startingBid, taxDeedPurchaser, winningBidAmount)
		}
	}

	err = db.InsertGwinnettPastSalesData(salesData, database)
	if err != nil {
		log.Fatalf("unable to insert data into db: %v", err)
	}
}
