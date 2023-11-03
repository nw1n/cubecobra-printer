package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCsvFile(fileName string) ([][]string, error) {
	// Open the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// remove header
	records = records[1:]

	return records, nil
}

func ReadDiffCsvFiles(cubeCsv, prevCsv string) ([][]string, error){
	cubeData, cErr := ReadCsvFile(cubeCsv)
	prevData, pErr := ReadCsvFile(prevCsv)

	
	if cErr != nil {
		return nil, cErr
	}
	if pErr != nil {
		return nil, pErr
	}

	newCards := make([][]string, 0)

	for _, cardRow := range cubeData {
		fmt.Println(cardRow[0])
		isInPrevCube := false
		for _, prevCardRow := range prevData {
				if prevCardRow[0] == cardRow[0] {
						isInPrevCube = true
						break
				}
		}
		if !isInPrevCube {
				newCards = append(newCards, cardRow)
		}
	}

	return newCards, nil
}

func ReadFullCsvData(CUBE_CSV_FILE, PREVIOUS_CUBE_CSV_FILE string) ([][]string){
	csvData := make([][]string, 0)
	isDiffMode := PREVIOUS_CUBE_CSV_FILE != ""

	if(isDiffMode){
		tmpCsvData, err := ReadDiffCsvFiles(CUBE_CSV_FILE, PREVIOUS_CUBE_CSV_FILE)
		if err != nil {
			panic(err)
		}
		csvData = tmpCsvData
	} else {
		tmpCsvData, err := ReadCsvFile(CUBE_CSV_FILE)
		if err != nil {
			panic(err)
		}
		csvData = tmpCsvData
	}
	return csvData
}