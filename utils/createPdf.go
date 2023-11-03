package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/jung-kurt/gofpdf"

	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"

	"golang.org/x/image/draw"
)

func copyFolder(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}
	if err := copyDir(src, dst); err != nil {
		return err
	}
	return nil
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destinationPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(sourcePath, destinationPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(sourcePath, destinationPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, sourceFile, 0644)
	if err != nil {
		return err
	}
	return nil
}

func trimBorders(borderedImagesFolder string, borderSize int) error {
	dir := borderedImagesFolder
	out := dir

	fillColor := color.RGBA{24, 21, 16, 255} // #181510 in RGB

	fmt.Println("folder:", dir)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".png") || strings.HasSuffix(info.Name(), ".jpg")) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			var img image.Image
			if strings.HasSuffix(info.Name(), ".png") {
				img, err = png.Decode(file)
			} else if strings.HasSuffix(info.Name(), ".jpg") {
				img, err = jpeg.Decode(file)
			}
			if err != nil {
				return err
			}

			originalBounds := img.Bounds()
			newBounds := image.Rect(0, 0, originalBounds.Dx()-2*borderSize, originalBounds.Dy()-2*borderSize)
			dst := image.NewRGBA(newBounds)

			// Fill the new image with the background color
			draw.Draw(dst, dst.Bounds(), &image.Uniform{fillColor}, image.ZP, draw.Src)

			// Copy the image content, excluding the original border
			draw.Draw(dst, newBounds, img, image.Pt(borderSize, borderSize), draw.Src)

			// Saving the image
			outFile, err := os.Create(filepath.Join(out, info.Name()))
			if err != nil {
				return err
			}
			defer outFile.Close()

			if strings.HasSuffix(info.Name(), ".png") {
				err = png.Encode(outFile, dst)
			} else if strings.HasSuffix(info.Name(), ".jpg") {
				err = jpeg.Encode(outFile, dst, nil)
			}
			if err != nil {
				return err
			}

			fmt.Println("Processed and saved image:", path)
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func addBorders(srcImagesFolder string, borderSize int) error {
	dir := srcImagesFolder
	out := dir

	borderColor := color.RGBA{24, 21, 16, 255} // Border color #181510 in RGB

	fmt.Println("folder:", dir)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".png") || strings.HasSuffix(info.Name(), ".jpg")) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			var img image.Image
			if strings.HasSuffix(info.Name(), ".png") {
				img, err = png.Decode(file)
			} else if strings.HasSuffix(info.Name(), ".jpg") {
				img, err = jpeg.Decode(file)
			}
			if err != nil {
				return err
			}

			originalBounds := img.Bounds()

			// Calculate new dimensions including the border
			newWidth := originalBounds.Dx() + 2*borderSize
			newHeight := originalBounds.Dy() + 2*borderSize

			// Create a new image with borders
			finalImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

			// Fill the entire image with the border color
			draw.Draw(finalImage, finalImage.Bounds(), &image.Uniform{borderColor}, image.ZP, draw.Src)

			// Draw the original image at the center (leaving space for borders)
			draw.Draw(finalImage, image.Rect(borderSize, borderSize, newWidth-borderSize, newHeight-borderSize), img, originalBounds.Min, draw.Src)

			// Save the image with borders
			outFile, err := os.Create(filepath.Join(out, info.Name()))
			if err != nil {
				return err
			}
			defer outFile.Close()

			if strings.HasSuffix(info.Name(), ".png") {
				err = png.Encode(outFile, finalImage)
			} else if strings.HasSuffix(info.Name(), ".jpg") {
				err = jpeg.Encode(outFile, finalImage, nil)
			}
			if err != nil {
				return err
			}

			fmt.Println("Processed and saved image:", path)
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func SortMapByValue(inputMap map[string]int) map[string]int {
	keys := make([]string, 0, len(inputMap))
	values := make([]int, 0, len(inputMap))

	for key, value := range inputMap {
		keys = append(keys, key)
		values = append(values, value)
	}

	sort.Slice(keys, func(i, j int) bool {
		return values[i] < values[j]
	})

	sortedMap := make(map[string]int)
	for _, key := range keys {
		sortedMap[key] = inputMap[key]
	}

	return sortedMap
}

func saveMapToJSON(inputMap map[string]int, filename string) error {
	println("Saving map to JSON file:", filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(inputMap)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}


func createPDF(directory, output string, cards []Card) error {
	dir, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	// if no images in folder, return
	if len(dir) == 0 {
		fmt.Printf("No images in folder %s. Skipping PDF creation for file %s.", directory, output)
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, entry := range dir {
		if !entry.IsDir() {
			fileName := entry.Name()

			imagePath := filepath.Join(directory, fileName)

			pdf.AddPage()
			pdf.Image(imagePath, 0, 0, 210, 297, false, "", 0, "")
		}
	}

	err = pdf.OutputFileAndClose(output)
	if err != nil {
		return err
	}

	return nil
}

func getLastElementAfterSlashSplit(input string) string {
	splitStr := strings.Split(input, "/")
	backSlashSplitStr := strings.Split(input, "\\")
	if len(backSlashSplitStr) > len(splitStr) {
		splitStr = backSlashSplitStr
	}
	if len(splitStr) > 0 {
		return splitStr[len(splitStr)-1]
	}
	return ""
}

func createDuplicates(borderedImagesFolder string, cards []Card) error {
	dir := borderedImagesFolder
	filesInDir, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, imageFile := range filesInDir {
		imageName := imageFile.Name()
		imagePath := filepath.Join(dir, imageName)
		
		imagePathWithoutExtension := strings.Replace(imagePath, ".png", "", 1)
		
		imageWithoutEnding := imagePathWithoutExtension
		imageWithoutEnding = strings.Replace(imageWithoutEnding, "_b", "", 1)
		imageWithoutEnding = strings.Replace(imageWithoutEnding, "_a", "", 1)

		isDuplicated := false
		amount := 1
		for _, card := range cards {
			cardImagePath := card.GetImagePath()
			cardImagePathWithoutExtension := strings.Replace(cardImagePath, ".png", "", 1)

			if getLastElementAfterSlashSplit(imageWithoutEnding) == getLastElementAfterSlashSplit(cardImagePathWithoutExtension) {
				isDuplicated = true
				amount = card.Amount
				break
			}
		}

		if isDuplicated {
			for i := 1; i < amount; i++ {
				numberStr := strconv.Itoa(i + 1)
				thePathRes := imageWithoutEnding
				thePathRes = thePathRes + "_" + numberStr

				if strings.HasSuffix(imagePath, "_a.png") {
					thePathRes = thePathRes + "_a"
				}

				if strings.HasSuffix(imagePath, "_b.png") {
					thePathRes = thePathRes + "_b"
				}
				thePathRes = thePathRes + ".png"

				// copy file
				err := copyFile(imagePath, thePathRes)
				if err != nil {
					PrintColorLn("ERROR!", "red")
					fmt.Println("Error copying file:", err)
					return err
				}
				fmt.Println("Copied file to: ", thePathRes)
			}
		}

	}

	return nil
}

func ProcessImagesToPdf(srcImagesFolder, borderedImagesFolder, finalPdfPath string, cards []Card) error {
	borderSize := 24

	err := copyFolder(srcImagesFolder, borderedImagesFolder)
	if err != nil {
		PrintColorLn("ERROR!", "red")
		fmt.Println("Error copying folder:", err)
		return err
	}

	err = trimBorders(borderedImagesFolder, borderSize)
	if err != nil {
		PrintColorLn("ERROR!", "red")
		fmt.Println("Error drawing borders:", err)
		return err
	}

	err = addBorders(borderedImagesFolder, borderSize * 2)
	if err != nil {
		PrintColorLn("ERROR!", "red")
		fmt.Println("Error drawing borders:", err)
		return err
	}

	err = createDuplicates(borderedImagesFolder, cards)
	if err != nil {
		PrintColorLn("ERROR!", "red")
		fmt.Println("Error creating duplicates:", err)
		return err
	}

	err = createPDF(borderedImagesFolder, finalPdfPath, cards)
	if err != nil {
		PrintColorLn("ERROR!", "red")
		fmt.Println("Error creating PDF:", err)
		return err
	}

	fmt.Println("")
	PrintColorLn("Finished PDF creation successfully", "green")
	fmt.Println("")
	return nil
}
