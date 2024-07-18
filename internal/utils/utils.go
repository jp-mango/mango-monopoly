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
	"strconv"
	"strings"
	"time"
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

func RunPdfExtraction(scriptPath string, county string) {
	fmt.Printf("Extracting data to csv located in ./tax-auction/%s. Please wait...\n", county)

	// create csv files
	if UpperTrim(county) == "GWINNETT" {
		os.Create(fmt.Sprintf("../../tax-auction/%s/csv/%s-Past-Sales.csv", county, county))
		os.Create(fmt.Sprintf("../../tax-auction/%s/csv/%s-Upcoming-Sales.csv", county, county))
	} else if UpperTrim(county) == "PAULDING" {
		os.Create(fmt.Sprintf("../../tax-auction/%s/csv/%s-Upcoming-Sales.csv", county, county))
	}

	cmd := exec.Command("cmd", "/C", ".venv\\Scripts\\activate && python", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		slog.Error(fmt.Sprintf("Error running Python script: %v", err))
	}

	fmt.Print("Data extraction complete\n\n")
}

func ReadCSV(filename string, county string) [][]string {
	dataPath := filepath.Join(fmt.Sprintf("tax-auction/%s/csv", county), filename)

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

func FormatMoney(value string) (float64, error) {
	cleaned := (strings.Replace(strings.Replace(value, "$", "", -1), ",", "", -1))

	return strconv.ParseFloat(cleaned, 64)
}

func StringToDate(date string) (time.Time, error) {
	t, err := time.Parse("January 2, 2006", date)
	if err != nil {
		return t, fmt.Errorf("unable to convert string to date")
	}

	return t, nil
}

func ASCIISpace(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, " ", "%20"))
}
