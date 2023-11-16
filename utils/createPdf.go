package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jung-kurt/gofpdf"

	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
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

func drawBlackTriangle(img draw.Image, colorVal color.Color, x1, y1, x2, y2 int) {
		// Set the fill color to black
		black := colorVal

		// Define the vertices of the triangle
		p1 := image.Point{x1, y1}
		p2 := image.Point{x2, y1}
		p3 := image.Point{(x1 + x2) / 2, y2}

		// Draw the triangle on the image
		draw.Draw(img, image.Rectangle{p1, p2}, &image.Uniform{black}, image.Point{}, draw.Over)
		draw.Draw(img, image.Rectangle{p2, p3}, &image.Uniform{black}, image.Point{}, draw.Over)
		draw.Draw(img, image.Rectangle{p3, p1}, &image.Uniform{black}, image.Point{}, draw.Over)
}

func saveImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
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

// TriangleMask creates a triangle mask with a specified base and height.
func TriangleMask(base, height int) *image.Alpha {
	mask := image.NewAlpha(image.Rect(0, 0, base, height))

	for y := 0; y < height; y++ {
		xEnd := base - y

		for x := 0; x < xEnd; x++ {
			mask.SetAlpha(x, y, color.Alpha{255})
		}
	}

	return mask
}

// rotate90 rotates the given mask 90 degrees counterclockwise.
func rotate90(mask *image.Alpha) *image.Alpha {
	bounds := mask.Bounds()
	rotated := image.NewAlpha(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rotated.SetAlpha(y, bounds.Max.X-x-1, mask.AlphaAt(x, y))
		}
	}

	return rotated
}

// rotate180 rotates the given mask 180 degrees counterclockwise.
func rotate180(mask *image.Alpha) *image.Alpha {
	return rotate90(rotate90(mask))
}

// rotate270 rotates the given mask 270 degrees counterclockwise.
func rotate270(mask *image.Alpha) *image.Alpha {
	return rotate90(rotate90(rotate90(mask)))
}



func addImageCorners(borderedImagesFolder string) error {
	dir := borderedImagesFolder
	out := dir

	// borderColor := color.RGBA{0, 0, 255, 255}
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
			newWidth := originalBounds.Dx()
			newHeight := originalBounds.Dy()

			// Create a new image with borders
			finalImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

			// draw the original image
			draw.Draw(finalImage, image.Rect(0, 0, newWidth, newHeight), img, originalBounds.Min, draw.Src)

			triangleBase := 31
			triangleHeight := 31
		
			// Create triangle masks for each corner
			topLeftMask := TriangleMask(triangleBase, triangleHeight)
			topRightMask := TriangleMask(triangleBase, triangleHeight)
			bottomLeftMask := TriangleMask(triangleBase, triangleHeight)
			bottomRightMask := TriangleMask(triangleBase, triangleHeight)
		
			// Rotate masks for top-right and bottom-right corners
			topRightMask = rotate90(topRightMask)
			bottomLeftMask = rotate270(bottomLeftMask)
			bottomRightMask = rotate180(bottomRightMask)
		
			draw.DrawMask(finalImage, finalImage.Bounds(), image.NewUniform(borderColor), image.Point{}, topLeftMask, image.Point{}, draw.Over)
			rotatedImage := imaging.Rotate90(finalImage)
			draw.DrawMask(rotatedImage, finalImage.Bounds(), image.NewUniform(borderColor), image.Point{}, topLeftMask, image.Point{}, draw.Over)
			rotatedImage = imaging.Rotate90(rotatedImage)
			draw.DrawMask(rotatedImage, finalImage.Bounds(), image.NewUniform(borderColor), image.Point{}, topLeftMask, image.Point{}, draw.Over)
			rotatedImage = imaging.Rotate90(rotatedImage)
			draw.DrawMask(rotatedImage, finalImage.Bounds(), image.NewUniform(borderColor), image.Point{}, topLeftMask, image.Point{}, draw.Over)
			rotatedImage = imaging.Rotate90(rotatedImage)

			finalImageX := rotatedImage	

			// Save the image with corners
			outFile, err := os.Create(filepath.Join(out, info.Name()))
			if err != nil {
				return err
			}
			defer outFile.Close()

			if strings.HasSuffix(info.Name(), ".png") {
				err = png.Encode(outFile, finalImageX)
			} else if strings.HasSuffix(info.Name(), ".jpg") {
				err = jpeg.Encode(outFile, finalImageX, nil)
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

	width := 280.00000
	height := 384.00000

	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, entry := range dir {
		if !entry.IsDir() {
			fileName := entry.Name()

			imagePath := filepath.Join(directory, fileName)

			pdf.AddPageFormat("P", gofpdf.SizeType{Wd: width, Ht: height})
			// Set the layout function for A4 page size
			pdf.SetAutoPageBreak(true, 0)
			pdf.SetMargins(0, 0, 0)
			pdf.SetFillColor(24, 21, 16) // Set RGB color for black
			pdf.Rect(0, 0, width, height, "F") // Rectangle covering the entire A4 page

			//pdf.Image(imagePath, 0, 0, 210, 297, false, "", 0, "")

			// Replace "your_image.jpg" with the path to your image file
			imageFile := imagePath

			// Embed the image in the PDF at its original dimensions
			pdf.Image(imageFile, 0, 0, width, height, false, "", 0, "")
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

	err = addImageCorners(borderedImagesFolder)
	if err != nil {
		PrintColorLn("ERROR!", "red")
		fmt.Println("Error drawing corners:", err)
		return err
	}

	err = addBorders(borderedImagesFolder, borderSize)
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
