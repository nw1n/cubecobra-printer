package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.com/bimbamdingdong/cubecobra-printer/utils"
)

const (
	// BASE
	TMP_FOLDER             = "cubecobra-printer-data"
	IMG_FOLDER             = "cubecobra-printer-data/downloaded-images"
	IMG_FLIP_FOLDER        = "cubecobra-printer-data/downloaded-images-flip"
	CSV_FOLDER             = "cubecobra-printer-data/csv"
	CUBE_CSV_FILE          = "cubecobra-printer-data/csv/cube-data.csv"
	PREVIOUS_CUBE_CSV_FILE = "cubecobra-printer-data/csv/previous-cube.csv"
	// PDF CREATION
	BORDERED_IMG_FOLDER      = "cubecobra-printer-data/bordered-images"
	BORDERED_FLIP_IMG_FOLDER = "cubecobra-printer-data/bordered-flip-images"
	PDF_FOLDER               = "cubecobra-printer-data/pdf"
	FINAL_PDF_PATH           = "cubecobra-printer-data/pdf/cube-one-sided.pdf"
	FINAL_FLIP_PDF_PATH      = "cubecobra-printer-data/pdf/cube-flip-cards.pdf"
	// TMP
	CSV_FLIP_EXAMPLE_URL      = "https://cubecobra.com/cube/download/csv/398e6e50-d585-43e0-ad97-6e4609da5e76"
	CSV_FLIP_DIFF_EXAMPLE_URL = "https://cubecobra.com/cube/download/csv/99898100-6f28-483a-9f25-556a848fe410"
	CSV_SIMPLE_EXAMPLE_URL    = "https://cubecobra.com/cube/download/csv/09f91955-989e-4abf-a470-b3763ba3b255"
	CSV_SIMPLE_NON_SINGLETON  = "https://cubecobra.com/cube/overview/simple-non-singleton"
	CSV_NON_SINGLETON_FLIP  = "https://cubecobra.com/cube/overview/non-singleton-flip"
	CSV_WIZARDS_CUBE_URL = "https://cubecobra.com/cube/list/3xfe3"
	CSV_MEGA_TEST_URL = "https://cubecobra.com/cube/list/mega-test"
	CSV_TINY_PEA_CUBE_URL = "https://cubecobra.com/cube/list/thetinypea" // singleton 180 cards
	CSV_REALISTIC_SIMPLE_URL = "https://cubecobra.com/cube/list/ed4473bb-54b7-455d-aa9b-d0b082b8047d" // singleton 20 cards

)

func runMainWithoutInput() {
	isPdfMode := true
	isDiffMode := false
	diffCsv := ""
	csvUrl := CSV_SIMPLE_EXAMPLE_URL
	diffCsvUrl := CSV_FLIP_DIFF_EXAMPLE_URL

	csvUrl = utils.GetCubeIdFromUrl(csvUrl)
	csvUrl = utils.GetCubecobraCsvUrl(csvUrl)

	utils.CreateAllTmpFolders([]string{TMP_FOLDER, IMG_FOLDER, IMG_FLIP_FOLDER, CSV_FOLDER, BORDERED_IMG_FOLDER, BORDERED_FLIP_IMG_FOLDER, PDF_FOLDER})

	utils.DownloadCsvFile(CUBE_CSV_FILE, csvUrl)
	if isDiffMode {
		diffCsv = PREVIOUS_CUBE_CSV_FILE
		diffCsvUrl = utils.GetCubeIdFromUrl(diffCsvUrl)
		diffCsvUrl = utils.GetCubecobraCsvUrl(diffCsvUrl)
		utils.DownloadCsvFile(PREVIOUS_CUBE_CSV_FILE, diffCsvUrl)
	}

	csvData := utils.ReadFullCsvData(CUBE_CSV_FILE, diffCsv)
	cards := utils.CreateAllCards(csvData, IMG_FOLDER, IMG_FLIP_FOLDER)
	utils.ProcessCardsAndDownloadImages(cards)

	isTestSuccess := utils.DoBasicChecks(cards, csvData, IMG_FOLDER)
	if !isTestSuccess {
		return
	}

	if !isPdfMode {
		return
	}

	// FLIP IMAGE PROCESSING
	utils.ProcessFlipImageForPpdf(IMG_FOLDER, IMG_FLIP_FOLDER)

	// CREATE PDF
	utils.ProcessImagesToPdf(IMG_FOLDER, BORDERED_IMG_FOLDER, FINAL_PDF_PATH, cards)
	utils.ProcessImagesToPdf(IMG_FLIP_FOLDER, BORDERED_FLIP_IMG_FOLDER, FINAL_FLIP_PDF_PATH, cards)
}

func runMainInteractive() {
	fmt.Println(utils.GetAsciiArtMainTitle())
	time.Sleep(200 * time.Millisecond)

	fmt.Println("Welcome to the Cubecobra Printer!")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)

	time.Sleep(500 * time.Millisecond)
	fmt.Println("This program will create a PDF of your cubecobra.com cube.")
	fmt.Println("")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("It will download all images in highest quality from scryfall.com, process them and create a PDF from it.")
	fmt.Println("This process will use a lot of Disk Space and can take a while.")
	fmt.Println("A 540 card cube will use up to 10 GB of disk space, so make sure you have enough disk space available.")
	fmt.Println("")
	fmt.Println("PLEASE NOTE: Use this program at your own risk.")
	fmt.Println("If you encounter issues, try deleting the created 'cubecobra-printer-data' folder and restart the process.")
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("This Process will create a new folder in your current directory called '%s'.\n", TMP_FOLDER)
	fmt.Println("")
	fmt.Println("If you want to continue type 'y' and press Enter. To exit enter 'n' or press Ctrl+C.")

	// Wait for input
	startInputStr, _ := reader.ReadString('\n')
	startInputStr = strings.TrimSpace(startInputStr)

	if startInputStr != "y" {
		fmt.Println("You did not type y. Aborting Process.")
		time.Sleep(1000 * time.Millisecond)
		return
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println("")
	fmt.Println("Initializing the cube creation process.")
	fmt.Println("")
	time.Sleep(200 * time.Millisecond)

	// check if folder already exists
	isTmpFolderExisting := utils.IsFolderExisting(TMP_FOLDER)
	if isTmpFolderExisting {
		fmt.Printf("The data folder already exists: %s\n", TMP_FOLDER)
		fmt.Println("The process will use the already downloaded data in this folder.")
		fmt.Printf("If you want to start from scratch, delete the folder '%s' manually and restart the process.\n", TMP_FOLDER)
		fmt.Println("IMPORTANT: If you are downloading a different cube than the one you downloaded before,")
		fmt.Println("you must delete the data folder manually, so that data will not be mixed up.")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("")
		fmt.Println("If you want to continue using the existing data type 'y' and press Enter. To exit enter 'n' or press Ctrl+C.")

		// Wait for input
		continueInputStr, _ := reader.ReadString('\n')
		continueInputStr = strings.TrimSpace(continueInputStr)

		fmt.Println("")
		if continueInputStr != "y" {
			fmt.Println("You did not type y. Aborting Process.")
			time.Sleep(1000 * time.Millisecond)
			return
		}
	} else {
		fmt.Println("Data folder does not exist yet.")
		fmt.Printf("The data folder will be created at this absolute path: %s\n", TMP_FOLDER)
	}

	cubeUrl := ""

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Please enter the URL of your cube's url:")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("(for example: https://cubecobra.com/cube/overview/09f91955-989e-4abf-a470-b3763ba3b255)")

	for {
		cubeUrl, _ = reader.ReadString('\n')
		cubeUrl = strings.TrimSpace(cubeUrl)

		if cubeUrl == "" {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("No valid input. Please enter a valid input for your cube's url.")
			continue
		}

		break
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println("")
	fmt.Println("If you want to only get the difference from previous cube enter the URL of your previous cube. (Experimental)")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("If you want to create a normal full cube just press Enter.")

	prevCubeUrl, _ := reader.ReadString('\n')
	prevCubeUrl = strings.TrimSpace(prevCubeUrl)

	cubeId := utils.GetCubeIdFromUrl(cubeUrl)
	prevCubeId := utils.GetCubeIdFromUrl(prevCubeUrl)

	time.Sleep(100 * time.Millisecond)
	if prevCubeId != "" {
		fmt.Printf("Your cube Id is: %s. You're prev cube id is: %s\n", cubeId, prevCubeId)
	} else {
		fmt.Printf("Your cube Id is: %s\n", cubeId)
	}
	fmt.Println("")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Your cube will be now be downloaded and processed.")
	fmt.Println("This process can take a while depending on the size of your cube.")
	fmt.Println("Please be patient and do not close the program.")
	fmt.Println("")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("If you want to continue type 'y' and press Enter. To exit enter 'n' or press Ctrl+C.")

	// Wait for input
	continueInputStr, _ := reader.ReadString('\n')
	continueInputStr = strings.TrimSpace(continueInputStr)

	fmt.Println("")
	if continueInputStr != "y" {
		fmt.Println("You did not type y. Aborting Process.")
		time.Sleep(1000 * time.Millisecond)
		return
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Start Downloading cube data...")
	fmt.Println("")
	time.Sleep(100 * time.Millisecond)

	isPdfMode := true
	// diff mode is dependent on prevCubeId
	isDiffMode := prevCubeId != ""
	diffCsvFilePath := ""
	csvUrl := utils.GetCubecobraCsvUrl(cubeId)
	diffCsvUrl := ""
	if isDiffMode {
		diffCsvUrl = utils.GetCubecobraCsvUrl(prevCubeId)
	}

	utils.CreateAllTmpFolders([]string{TMP_FOLDER, IMG_FOLDER, IMG_FLIP_FOLDER, CSV_FOLDER, BORDERED_IMG_FOLDER, BORDERED_FLIP_IMG_FOLDER, PDF_FOLDER})

	fmt.Printf("Downloading cube data from: %s\n", csvUrl)
	utils.DownloadCsvFile(CUBE_CSV_FILE, csvUrl)
	if isDiffMode {
		diffCsvFilePath = PREVIOUS_CUBE_CSV_FILE
		utils.DownloadCsvFile(PREVIOUS_CUBE_CSV_FILE, diffCsvUrl)
	}

	csvData := utils.ReadFullCsvData(CUBE_CSV_FILE, diffCsvFilePath)
	cards := utils.CreateAllCards(csvData, IMG_FOLDER, IMG_FLIP_FOLDER)
	utils.ProcessCardsAndDownloadImages(cards)

	//isTestSuccess := utils.DoBasicChecks(cards, csvData, IMG_FOLDER)
	//if !isTestSuccess {
	//	return
	//}

	if !isPdfMode {
		return
	}

	// FLIP IMAGE PROCESSING
	utils.ProcessFlipImageForPpdf(IMG_FOLDER, IMG_FLIP_FOLDER)

	// CREATE PDF
	utils.ProcessImagesToPdf(IMG_FOLDER, BORDERED_IMG_FOLDER, FINAL_PDF_PATH, cards)
	utils.ProcessImagesToPdf(IMG_FLIP_FOLDER, BORDERED_FLIP_IMG_FOLDER, FINAL_FLIP_PDF_PATH, cards)

	fmt.Println("")
	fmt.Println("----------------------------------------")
	fmt.Println("")
	utils.PrintColorLn("SUCCESS!", "green")
	fmt.Println("")
	fmt.Println("Finished creating PDFs.")
	fmt.Println("")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("The PDFs are located in the following folder:")
	fmt.Println(TMP_FOLDER + "/pdf")
	// press any key to exit
	time.Sleep(200 * time.Millisecond)
	fmt.Println("")
	fmt.Println("Press enter to exit the program.")
	reader.ReadString('\n')
}

func main() {
	// Define a "dev" flag of type boolean
	devPtr := flag.Bool("dev", false, "Set development mode")
	// Parse the command-line flags
	flag.Parse()
	// set dev mode from flag
	isDevMode := *devPtr

	if isDevMode {
		runMainWithoutInput()
	} else {
		runMainInteractive()
	}
}
