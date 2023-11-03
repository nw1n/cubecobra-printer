package utils

import (
	"fmt"
	"os"
	"strings"
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
  ______             __                                              
 /      \           /  |                                             
/$$$$$$  | __    __ $$ |____    ______                               
$$ |  $$/ /  |  /  |$$      \  /      \                              
$$ |      $$ |  $$ |$$$$$$$  |/$$$$$$  |                             
$$ |   __ $$ |  $$ |$$ |  $$ |$$    $$ |                             
$$ \__/  |$$ \__$$ |$$ |__$$ |$$$$$$$$/                              
$$    $$/ $$    $$/ $$    $$/ $$       |                             
 $$$$$$/   $$$$$$/  $$$$$$$/   $$$$$$$/                              
  ______                                  __                         
 /      \                                /  |                        
/$$$$$$  |  ______   ______    ______   _$$ |_     ______    ______  
$$ |  $$/  /      \ /      \  /      \ / $$   |   /      \  /      \ 
$$ |      /$$$$$$  /$$$$$$  | $$$$$$  |$$$$$$/   /$$$$$$  |/$$$$$$  |
$$ |   __ $$ |  $$/$$    $$ | /    $$ |  $$ | __ $$ |  $$ |$$ |  $$/ 
$$ \__/  |$$ |     $$$$$$$$/ /$$$$$$$ |  $$ |/  |$$ \__$$ |$$ |      
$$    $$/ $$ |     $$       |$$    $$ |  $$  $$/ $$    $$/ $$ |      
 $$$$$$/  $$/       $$$$$$$/  $$$$$$$/    $$$$/   $$$$$$/  $$/                                                

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

func PrintColorLn(msg string, color string) {
	reset := "\033[0m"
	colorCode := reset

	cRed := "\033[31m"
	cGreen := "\033[32m"
	cYellow := "\033[33m"

	if color == "red" {
		colorCode = cRed
	} else if color == "green" {
		colorCode = cGreen
	} else if color == "yellow" {
		colorCode = cYellow
	}



	fmt.Println(colorCode + msg + reset)
}