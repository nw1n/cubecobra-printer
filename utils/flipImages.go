package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ProcessFlipImageForPpdf(IMG_FOLDER string, IMG_FLIP_FOLDER string) {
	// EXTRA FLIP PROCESSING
	files, _ := os.ReadDir(IMG_FLIP_FOLDER)
	for _, file := range files {
		filename := file.Name()
		isContainingB := strings.Contains(filename, "_b.png")
		isContainingA := strings.Contains(filename, "_a.png")
		if !isContainingB && !isContainingA {
			if strings.HasSuffix(filename, ".png") {
				oldPath := filepath.Join(IMG_FLIP_FOLDER, filename)
				newPath := filepath.Join(IMG_FLIP_FOLDER, strings.Replace(filename, ".png", "_b.png", 1))
				os.Rename(oldPath, newPath)
			}
		}
		shortFileName := strings.Replace(filename, "_b.png", ".png", 1)
		srcPath := filepath.Join(IMG_FOLDER, shortFileName)
		if _, err := os.Stat(srcPath); err == nil {
			newFileName := strings.Replace(shortFileName, ".png", "_a.png", 1)
			destPath := filepath.Join(IMG_FLIP_FOLDER, newFileName)
			os.Rename(srcPath, destPath)
		}
		//fmt.Println(shortFileName)
	}
}
