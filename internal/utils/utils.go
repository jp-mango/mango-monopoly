package utils

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
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

func RunPdfExtraction(script string) {
	fmt.Println("Extracting data to csv located in ./tax-auction/Gwinnett. Please wait...")

	// create csv files
	os.Create("../../tax-auction/Gwinnett/csv/Gwinnett-Past-Sales.csv")
	os.Create("../../tax-auction/Gwinnett/csv/Gwinnett-Upcoming-Sales.csv")

	cmd := exec.Command("cmd", "/C", ".venv\\Scripts\\activate && python", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		slog.Error(fmt.Sprintf("Error running Python script: %v", err))
	}

	fmt.Print("Data extraction complete\n\n")
}
