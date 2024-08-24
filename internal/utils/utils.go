package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
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

func GetLocationInfo(addr string) []string {
	type AddrInfo struct {
		Addy string `json:"display_name"`
		Lat  string `json:"lat"`
		Lon  string `json:"lon"`
	}

	geocodeURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", url.QueryEscape(addr))

	client := &http.Client{}
	req, err := http.NewRequest("GET", geocodeURL, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return nil
	}
	req.Header.Set("User-Agent", "mango-monopoly/1.0")

	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
	}

	fmt.Print(string(body))

	var result []AddrInfo
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error parsing JSON response: %v", err)
	}

	var address, latitude, longitude string

	if len(result) > 0 {
		latitude = result[0].Lat
		longitude = result[0].Lon
		address = result[0].Addy
	}

	return []string{latitude, longitude, address}
}
