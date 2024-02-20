package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

)

type Card struct {
    CsvCardData []string
    Title string
    NameOfSet  string
    CollectorNumber string
    ImageBasePath string
    FlipImageBasePath string
		ScryfallData map[string]interface{}
    ScryfallImgUri string
		Layout string
		ScryFallImgUriFlip string
		ScryFallCardFaces map[string]interface{}
		Amount int
}

func (self *Card) FetchScryfallData() {
	resp, err := http.Get(self.GetCardQueryUri()) // Assuming getCardQueryUri() is defined
	if err != nil {
		fmt.Println("Error fetching Scryfall Data for:", self.Title)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", resp.StatusCode)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&self.ScryfallData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	time.Sleep(50 * time.Millisecond)
}

func (c *Card) DownloadImage() {
	uri := c.GetScryfallImgUri()
	if uri == "" {
			fmt.Printf("Scryfall data missing. can not download image: %s\n", c.Title)
			return
	}

	basePath := c.GetImageBasePath()
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
			os.MkdirAll(basePath, os.ModePerm)
	}

	outPath := c.GetImagePath()
	out, err := os.Create(outPath)
	if err != nil {
			fmt.Printf("Error creating file %s: %v\n", outPath, err)
			return
	}
	defer out.Close()

	resp, err := http.Get(uri)
	if err != nil {
			fmt.Printf("Error downloading image for %s: %v\n", c.Title, err)
			return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", outPath, err)
			return
	}

	fmt.Printf("Downloaded: %s\n", outPath)
	time.Sleep(50 * time.Millisecond)
}

func (c *Card) DownloadFlipImage() {
	if !c.IsFlipCard() {
			return
	}

	uri := c.GetScryfallImgUriFlip()
	if uri == "" {
			fmt.Printf("Scryfall data missing. can not download flip image: %s\n", c.Title)
			return
	}

	basePath := c.GetFlipImageBasePath()
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
			os.MkdirAll(basePath, os.ModePerm)
	}

	outPath := c.GetFlipImagePath()
	out, err := os.Create(outPath)
	if err != nil {
			fmt.Printf("Error creating file %s: %v\n", outPath, err)
			return
	}
	defer out.Close()

	resp, err := http.Get(uri)
	if err != nil {
			fmt.Printf("Error downloading flip image for %s: %v\n", c.Title, err)
			return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", outPath, err)
			return
	}

	fmt.Printf("Downloaded: %s\n", outPath)
	time.Sleep(50 * time.Millisecond)
}

func (card *Card) IsFlipCard() bool {
	if card.ScryfallData["layout"] == "adventure" {
			return false
	}
	if card.ScryfallData["layout"] == "split" {
			return false
	}
	if _, ok := card.ScryfallData["card_faces"]; ok {
			return true
	}
	return false
}

func (card *Card) InitScryFallData() {
    // card.ScryfallImgUri = card.ScryfallData["image_uris"].(map[string]interface{})["png"].(string)
		card.Layout = card.ScryfallData["layout"].(string)
}

func (c *Card) GetCardQueryUri() string {
	return fmt.Sprintf("https://api.scryfall.com/cards/%s/%s", c.NameOfSet, c.CollectorNumber)
}

func (c *Card) GetImageBasePath() string {
	return c.ImageBasePath
}

func (c *Card) GetFlipImageBasePath() string {
	return c.FlipImageBasePath
}

func (c *Card) GetScryfallImgUriFlip() string {
	if c.IsFlipCard() {
			return c.ScryfallData["card_faces"].([]interface{})[1].(map[string]interface{})["image_uris"].(map[string]interface{})["png"].(string)
	}
	return ""
}

func (c *Card) GetScryfallCardUri() string {
	if _, ok := c.ScryfallData["scryfall_uri"]; ok {
			return c.ScryfallData["scryfall_uri"].(string)
	}
	return ""
}

func (c *Card) GetScryfallImgUri() string {
	if c.IsFlipCard() {
			return c.ScryfallData["card_faces"].([]interface{})[0].(map[string]interface{})["image_uris"].(map[string]interface{})["png"].(string)
	} else {
			if _, ok := c.ScryfallData["image_uris"]; ok {
					return c.ScryfallData["image_uris"].(map[string]interface{})["png"].(string)
			}
	}
	return ""
}

func (c *Card) GetImageFileName() string {
	return c.GetHyphenatedName() + "__" + c.GetId() + ".png"
}

func (c *Card) GetImagePath() string {
	return filepath.Join(c.GetImageBasePath(), c.GetImageFileName())
}



func (c *Card) GetFlipImagePath() string {
	return filepath.Join(c.GetFlipImageBasePath(), c.GetImageFileName())
}

func (c *Card) IsFileOnDisk() bool {
	_, err := os.Stat(c.GetImagePath())
	isFileExisting := !os.IsNotExist(err)
	return isFileExisting
}

func (c *Card) IsFlipFileOnDisk() bool {
	_, err := os.Stat(c.GetFlipImagePath())
	return !os.IsNotExist(err)
}

func (c *Card) GetId() string {
	idStr := RemoveNonASCIIChars(c.NameOfSet) + "-" +  RemoveNonASCIIChars(c.CollectorNumber)
	// log the idStr
	fmt.Println(idStr)
	return idStr
}

func (c *Card) GetHyphenatedName() string {
	res := strings.ToLower(c.Title)
	res = strings.ReplaceAll(res, "'", "")
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	res = reg.ReplaceAllString(res, "-")
	return res
}

func NewCard(csvCardData []string, imageBasePath string, flipImageBasePath string) Card {
    return Card{
        CsvCardData: csvCardData,
        ImageBasePath: imageBasePath,
        FlipImageBasePath: flipImageBasePath,
    }
}

func (card *Card) Init() {
    card.Title = card.CsvCardData[0]
    card.NameOfSet = card.CsvCardData[4]
    card.CollectorNumber = card.CsvCardData[5]
    card.CsvCardData = nil
		card.Amount = 1
}

func CreateAllCards(csvData [][]string, imageBasePath string, flipImageBasePath string) []Card {
    var cards []Card // create an empty slice of Card pointers
    
    // loop through the CSV data and create a Card for each row
    for _, row := range csvData {
        card := NewCard(row, imageBasePath, flipImageBasePath)
        card.Init()
        cards = append(cards, card)
    }

		// make cards unique
		var cardsUnique []Card
		var cardsUniqueMap = make(map[string]bool)
		for _, card := range cards {
			if _, value := cardsUniqueMap[card.Title]; !value {
				cardsUniqueMap[card.Title] = true
				cardsUnique = append(cardsUnique, card)
			} else {
				for i, cardUnique := range cardsUnique {
					if cardUnique.Title == card.Title {
						cardsUnique[i].Amount++
					}
				}
			}
		}

    return cardsUnique
}

func ProcessCardsAndDownloadImages(cards []Card) []Card {
	for _, card := range cards {
		if card.IsFileOnDisk() {
			// if false {
			fmt.Println("File already exists:", card.GetImagePath())
		} else {
			card.FetchScryfallData()
			card.InitScryFallData()
			card.DownloadImage()
			card.DownloadFlipImage()
		}
	}
	return cards
}

func getFileSize(filename string) (int64, error) {
	file, err := os.Open(filename)
	if err != nil {
			return 0, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
			return 0, err
	}

	// Get the file size in bytes
	fileSize := fileInfo.Size()

	return fileSize, nil
}
