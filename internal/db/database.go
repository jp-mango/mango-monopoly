package db

import (
	"database/sql"
	"fmt"
	"log"
	"mango-monopoly/internal/scraper"
	"mango-monopoly/internal/utils"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func DbConnect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_CON")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to mango-monopoly")

	return db, nil
}

func ResetDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_CON")
	if dsn == "" {
		log.Fatal("DB_CON environment variable is not set")
	}

	fmt.Printf("Resetting database. Please wait...\n\n")

	commands := []string{
		fmt.Sprintf("migrate -path ./migrations -database %s force 1", dsn),
		fmt.Sprintf("migrate -path ./migrations -database %s down -all", dsn),
		fmt.Sprintf("migrate -path ./migrations -database %s up", dsn),
	}

	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/c"
	} else {
		shell = "sh"
		flag = "-c"
	}

	for _, cmdStr := range commands {
		cmd := exec.Command(shell, flag, cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Fatalf("Error resetting db: %v", err)
		}
	}

	fmt.Println("\nDatabase reset.")
	fmt.Println()
}

func LoadCounties(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO counties (name, state)
		SELECT 'Gwinnett', 'Georgia'
		WHERE
		NOT EXISTS(
		SELECT name from counties WHERE name = 'Gwinnett'
		);
		`)
	if err != nil {
		return fmt.Errorf("error loading counties: %v", err)
	}

	fmt.Println("Counties loaded into mango-monopoly")

	return nil
}

func InsertGwinnettPastSalesData(salesData [][]string, db *sql.DB) (rowsEffected int64, err error) {
	header := salesData[0]
	phrase1 := "TENTATIVELY SCHDULED TAX SALE"
	phrase2 := "WE DID NOT HAVE A DECEMBER 5, 2023 TAX SALE"
	phrase3 := "NO BID"
	var totalRowsAffected int64

	query := `
		INSERT INTO Past_Sales (auction_date, parcel_id, previous_owner, addr, starting_bid, tax_deed_purchaser, winning_bid_amount)
		SELECT $1, CAST ($2 AS VARCHAR), $3, $4, $5, $6, $7
		WHERE NOT EXISTS(
		SELECT 1 FROM Past_Sales WHERE parcel_id = $2
		);`

	for i, value := range salesData {
		if slices.Compare(value, header) == 0 {
			// Skip header row
			continue
		}

		auctionDate := utils.UpperTrim(value[0])
		parcelID := utils.UpperTrim(value[1])
		previousOwner := utils.UpperTrim(value[2])
		addr := utils.UpperTrim(value[3])
		startingBidStr := utils.UpperTrim(value[5])
		taxDeedPurchaser := utils.UpperTrim(value[6])
		winningBidAmountStr := utils.UpperTrim(value[7])

		// Skip rows based on specific phrases for previousOwner
		if previousOwner == phrase1 || previousOwner == phrase2 {
			continue
		}

		var startingBid, winningBidAmount float64
		var err error

		// Handle "NO BID" separately for startingBid
		if startingBidStr == phrase3 {
			startingBid = 0
		} else {
			startingBid, err = utils.FormatMoney(startingBidStr)
			if err != nil {
				return 0, fmt.Errorf("error converting string to float at index %d: %v", i, err)
			}
		}

		// Handle "NO BID" separately for winningBidAmount
		if winningBidAmountStr == phrase3 {
			winningBidAmount = 0
		} else {
			winningBidAmount, err = utils.FormatMoney(winningBidAmountStr)
			if err != nil {
				return 0, fmt.Errorf("error converting string to float at index %d: %v", i, err)
			}
		}
		result, err := db.Exec(query, auctionDate, parcelID, previousOwner, addr, startingBid, taxDeedPurchaser, winningBidAmount)
		if err != nil {
			return 0, fmt.Errorf("error inserting data at index %d: %v", i, err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return totalRowsAffected, fmt.Errorf("error retrieving rows affected at index %d: %v", i, err)
		}
		totalRowsAffected += rowsAffected
	}

	fmt.Println("Successfully inserted gwinnett past auction data to db.")
	return totalRowsAffected, nil
}

func InsertGwinnettUpcomingSalesData(salesData [][]string, db *sql.DB) (rowsEffected int64, err error) {
	var totalRowsAffected int64
	auction_date, err := utils.StringToDate(scraper.GwinnettUpcomingAuctions[0])
	if err != nil {
		return 0, fmt.Errorf("error converting string to date: %v", err)
	}

	query := `
		INSERT INTO Upcoming_Sales (parcel_id, owner, auction_date,address, amount_due)
		SELECT CAST ($1 AS VARCHAR), $2, $3, $4, $5
		WHERE NOT EXISTS(
		SELECT 1 FROM Upcoming_Sales WHERE parcel_id = $1
		);`

	for i, value := range salesData {
		if i == 0 {
			continue
		}

		parcelID := utils.UpperTrim(value[0])
		owner := utils.UpperTrim(value[1])
		address := utils.UpperTrim(value[2])
		owed, err := utils.FormatMoney(utils.UpperTrim(value[3]))
		if err != nil {
			return 0, fmt.Errorf("error formatting string value at index %d: %v", i, err)
		}

		result, err := db.Exec(query, parcelID, owner, auction_date, address, owed)
		if err != nil {
			return 0, fmt.Errorf("error inserting at index %d: %v", i, err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return totalRowsAffected, fmt.Errorf("error retrieving rows affected at index %d: %v", i, err)
		}
		totalRowsAffected += rowsAffected
	}

	fmt.Println("Successfully inserted gwinnett upcoming auction data to db.")
	return totalRowsAffected, nil
}

func gwinnettTaxWorker(parcelIDs <-chan string, results chan<- int64, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	for parcel_id := range parcelIDs {
		details := scraper.TaxAssessorsOfficePull(parcel_id)

		updateQuery := `
			UPDATE Properties
			SET situs = $1,
				county_id = 1,
				property_type = $2,
				land_value = $3,
				building_value = $4,
				fair_market_value = $5,
				lot_size = $6,
				tax_assessor_url = $7
			WHERE parcel_id = $8
			AND tax_assessor_url IS NULL;
		`
		landValue, err := utils.FormatMoney(details.LandValue)
		if err != nil {
			log.Println("error formatting land_value", "value", details.LandValue, "error", err)
			results <- 0
			continue
		}

		buildingValue, err := utils.FormatMoney(details.BuildingValue)
		if err != nil {
			log.Println("error formatting building_value", "value", details.BuildingValue, "error", err)
			results <- 0
			continue
		}

		fairMarketValue, err := utils.FormatMoney(details.FairMarketValue)
		if err != nil {
			log.Println("error formatting fair_market_value", "value", details.FairMarketValue, "error", err)
			results <- 0
			continue
		}

		_, err = db.Exec(updateQuery, details.Address, details.PropertyType, landValue, buildingValue, fairMarketValue, details.LotSize, details.TaxAssessorURL, parcel_id)
		if err != nil {
			log.Println("error updating row", "parcel_id", parcel_id, "error", err)
			results <- 0
			continue
		}
		results <- details.Updates
	}
}

func UpdatePropertiesTable_Tax(db *sql.DB) (int64, error) {
	var totalRequests int64
	query := `SELECT parcel_id FROM Properties WHERE tax_assessor_url IS NULL;`
	var parcel_id string

	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	parcelIDs := make(chan string, 100)
	results := make(chan int64, 100)
	var wg sync.WaitGroup

	numWorkers := 20
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go gwinnettTaxWorker(parcelIDs, results, db, &wg)
	}

	fmt.Println("Pulling information from Gwinnett tax assessor, please wait... ")

	go func() {
		for rows.Next() {
			if err := rows.Scan(&parcel_id); err != nil {
				log.Println("error scanning row", "error", err)
				continue
			}
			parcelIDs <- parcel_id
		}
		close(parcelIDs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for updates := range results {
		totalRequests += updates
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error during rows iteration: %v", err)
	}

	return totalRequests, nil
}
