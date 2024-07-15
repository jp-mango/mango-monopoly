package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

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

func InsertGwinnettPastSalesData(salesData [][]string, db *sql.DB) error {
	header := salesData[0]
	phrase := "tentatively schduled tax sale"

	query := `
		INSERT INTO Past_Sales (auction_date, parcel_id, previous_owner, addr, starting_bid, tax_deed_purchaser, winning_bid_amount)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS(
		SELECT parcel_id FROM Past_Sales WHERE parcel_id = $2
		);`

	for i, value := range salesData {
		if slices.Compare(value, header) != 0 && strings.ToLower(strings.TrimSpace(value[2])) != phrase {
			auctionDate := value[0]
			parcelID := value[1]
			previousOwner := value[2]
			addr := value[3]
			startingBid := value[5]
			taxDeedPurchaser := value[6]
			winningBidAmount := value[7]

			_, err := db.Exec(query, auctionDate, parcelID, previousOwner, addr, startingBid, taxDeedPurchaser, winningBidAmount)
			if err != nil {
				return fmt.Errorf("error inserting data at index %d: %v", i, err)
			}
		}
	}

	fmt.Println("Successfully inserted gwinnett past auction data to db.")
	return nil
}
