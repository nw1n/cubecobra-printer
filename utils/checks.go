package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DoBasicChecks(cards []Card, csvData [][]string, IMG_FOLDER string) bool {
	// test cards
	fmt.Println("")
	fmt.Println("------------------------------------------")
	fmt.Println("Running tests")
	testSuccess := RunCardChecks(cards, csvData, IMG_FOLDER)

	fmt.Println("")
	if testSuccess {
		fmt.Println("------------------------------------------")
		PrintColorLn("All tests successful!", "green")
		fmt.Println("")
		return true
	}
	fmt.Println("------------------------------------------")
	PrintColorLn("ERROR: Some tests failed!", "red")
	fmt.Println("Please try running the script again.")
	fmt.Println("If the error persists, try deleting the 'cubecobra-printer-data' folder and restart the process.")
	fmt.Println("")

	return false
}

func RunCardChecks(cards []Card, csvData [][]string, imagesFolder string) bool {
	imageFiles, _ := os.ReadDir(imagesFolder)

	// uniqueRowsCsv
	uniqueRowsCsv := getCsvWithUniqueRows(csvData)
	csvData = uniqueRowsCsv
	// count number of cards in csv
	amountOfCardsInCsv := len(csvData)
	// count number of files in folder
	amountOfFiles := len(imageFiles)

	// check if number of cards in csv matches number of files in folder
	fmt.Printf("  Total Cards in List: %d\n", amountOfCardsInCsv)
	fmt.Printf("  Total Image Files in Folder: %d\n", amountOfFiles)
	if amountOfCardsInCsv == amountOfFiles {
		fmt.Println("Test Success. Number of unique Cards in List matches number of Files in Folder.")
	} else {
		fmt.Println("Number of unique Cards in List does not Match number of Files in Folder.")
		return false
	}

	// check file integrity of all images
	for _, imageFile := range imageFiles {
		imagePath := filepath.Join(imagesFolder, imageFile.Name())
		// check image file Size
		fileSize, _ := getFileSize(imagePath)

		// check if is folder
		if imageFile.IsDir() {
			fmt.Println("\nError: Image is a folder:", imagePath)
			return false
		}

		// check if file is png
		if !strings.HasSuffix(imageFile.Name(), ".png") {
			fmt.Println("\nError: Image is not a png: ", imagePath)
			return false
		}

		// check if file is at least 0.1 MB
		if(fileSize < 100000) {
			fmt.Println("\nFilesize too small:", imagePath)
			return false
		}
	}
	fmt.Println("Filesize check successful")


	for _, card := range cards {
		if card.IsFileOnDisk() {
			fmt.Println("Test Success. Image exists:", card.GetImagePath())
		} else {
			fmt.Println("File does not exist:", card.GetImagePath())
			return false
		}
		if card.IsFlipCard() {
			fmt.Println("Is Flip Card:", card.Title)
			if card.IsFlipFileOnDisk() {
				fmt.Println("Test Success. Flip Image exists:", card.GetFlipImagePath())
			} else {
				fmt.Println("File does not exist:", card.GetFlipImagePath())
				return false
			}
		}
	}
	return true
}
