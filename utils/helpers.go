package utils

import (
	"fmt"
	"os"
	"strings"
	"github.com/fatih/color"
)

func CreateFolderIfNotExisting(folderPath string) error {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderPath, 0755) // 0755 is the permission mode for the directory
		if errDir != nil {
			return errDir
		}
		fmt.Printf("Directory created: %s\n", folderPath)
	} else if err != nil {
		return err
	}
	return nil
}

func IsFolderExisting(folderPath string) bool {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateAllTmpFolders(sliceOfFolders []string) {
	for _, folder := range sliceOfFolders {
		CreateFolderIfNotExisting(folder)
	}
}

func GetCubeIdFromUrl(url string) string {
	url = strings.TrimSpace(url)
	splitUrl := strings.Split(url, "/")
	cubeId := splitUrl[len(splitUrl)-1]
	return cubeId
}

func GetCubecobraCsvUrl(inputStr string) string {
	cubeId := GetCubeIdFromUrl(inputStr)
	return fmt.Sprintf("https://cubecobra.com/cube/download/csv/%s", cubeId)
}

func GetAsciiArtMainTitle() string {
	return `

██████╗██╗   ██╗██████╗ ███████╗ ██████╗ ██████╗ ██████╗ ██████╗  █████╗ 
██╔════╝██║   ██║██╔══██╗██╔════╝██╔════╝██╔═══██╗██╔══██╗██╔══██╗██╔══██╗
██║     ██║   ██║██████╔╝█████╗  ██║     ██║   ██║██████╔╝██████╔╝███████║
██║     ██║   ██║██╔══██╗██╔══╝  ██║     ██║   ██║██╔══██╗██╔══██╗██╔══██║
╚██████╗╚██████╔╝██████╔╝███████╗╚██████╗╚██████╔╝██████╔╝██║  ██║██║  ██║
 ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝
																																					
██████╗ ██████╗ ██╗███╗   ██╗████████╗███████╗██████╗                     
██╔══██╗██╔══██╗██║████╗  ██║╚══██╔══╝██╔════╝██╔══██╗                    
██████╔╝██████╔╝██║██╔██╗ ██║   ██║   █████╗  ██████╔╝                    
██╔═══╝ ██╔══██╗██║██║╚██╗██║   ██║   ██╔══╝  ██╔══██╗                    
██║     ██║  ██║██║██║ ╚████║   ██║   ███████╗██║  ██║                    
╚═╝     ╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝                    
																																																												
`
}

func getCsvWithUniqueRows(csvData[][]string) [][]string {
	var uniqueRows [][]string
	var uniqueRowsMap = make(map[string]bool)
	for _, row := range csvData {
		if _, value := uniqueRowsMap[row[0]]; !value {
			uniqueRowsMap[row[0]] = true
			uniqueRows = append(uniqueRows, row)
		}
	}
	return uniqueRows
}

func PrintColorLn(msg string, colorStr string) {
	if colorStr == "red" {
		color.Red(msg)
	} else if colorStr == "green" {
		color.Green(msg)
	} else if colorStr == "yellow" {
		color.Yellow(msg)
	}
}
