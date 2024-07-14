package main

import (
	"mango-monopoly/internal/db"
	"mango-monopoly/internal/scraper"
)

func main() {
	scraper.DownloadGwinnettAuctionData()
	db.DbConnect()
}
