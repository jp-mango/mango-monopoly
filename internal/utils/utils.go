package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Function to download a file from a URL
func DownloadFile(filepath string, url string) error {
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

func RunPdfExtraction(scriptPath string) {
	fmt.Println("Extracting data to csv located in ./tax-auction/Gwinnett. Please wait...")

	// create csv files
	os.Create("../../tax-auction/Gwinnett/csv/Gwinnett-Past-Sales.csv")
	os.Create("../../tax-auction/Gwinnett/csv/Gwinnett-Upcoming-Sales.csv")

	cmd := exec.Command("cmd", "/C", ".venv\\Scripts\\activate && python", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		slog.Error(fmt.Sprintf("Error running Python script: %v", err))
	}

	fmt.Print("Data extraction complete\n\n")
}

func ReadCSV(filename string) [][]string {
	dataPath := filepath.Join("tax-auction/Gwinnett/csv", filename)

	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatalf("Unable to open CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse CSV: %v", err)
	}

	return records
}

func UpperTrim(word string) string {
	return strings.ToUpper(strings.TrimSpace(word))
}
