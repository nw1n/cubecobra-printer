package utils

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
)


func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func DownloadCsvFile(localFilePath string, csvURL string)  {

	//if fileExists(localFilePath) {
	//	fmt.Println("File already exists")
	//	return
	//}

	// Make an HTTP GET request to fetch the CSV data
	response, err := http.Get(csvURL)
	if err != nil {
		fmt.Println("Error fetching CSV:", err)
		return
	}
	defer response.Body.Close()

	// Parse the HTTP response
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: HTTP Status", response.StatusCode)
		return
	}

	// Create a CSV reader to read the downloaded data
	reader := csv.NewReader(response.Body)

	// Read and process the CSV data
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	file, err := os.Create(localFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	fmt.Println("Creating Local CSV File")

	writer := csv.NewWriter(file)
	writer.WriteAll(records)
	writer.Flush()
}
