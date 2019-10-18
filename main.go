package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

func main() {
	var inputPath string
	flag.StringVar(&inputPath, "i", ".", "Input path to folder")
	var outputPath string
	flag.StringVar(&outputPath, "o", "output.pdf", "Output path for pdf")
	flag.Parse()

	folders, files := scanRecursive(inputPath, []string{})
	fmt.Println("Found " + strconv.Itoa(len(folders)) + " folders containing " + strconv.Itoa(len(files)) + " files.")

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for _, file := range files {
		pdf.AddPage()

		imageW, imageH := getImageDimension(file)
		scale := math.Min(gopdf.PageSizeA4.W/float64(imageW), gopdf.PageSizeA4.H/float64(imageH))
		var rect gopdf.Rect
		rect.W = float64(imageW) * scale
		rect.H = float64(imageH) * scale

		centerW := (gopdf.PageSizeA4.W - rect.W) / 2
		centerH := (gopdf.PageSizeA4.H - rect.H) / 2

		err := pdf.Image(file, centerW, centerH, &rect)
		if err != nil {
			log.Panic(err)
		}
	}

	pdf.WritePdf(outputPath)
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Panic(err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Panic(file.Name(), err)
	}
	return image.Width, image.Height
}

func scanRecursive(dirPath string, ignore []string) ([]string, []string) {

	folders := []string{}
	files := []string{}

	// Scan
	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {

		_continue := false

		// Loop : Ignore Files & Folders
		for _, i := range ignore {

			// If ignored path
			if strings.Index(path, i) != -1 {

				// Continue
				_continue = true
			}
		}

		if _continue == false {

			f, err = os.Stat(path)

			// If no error
			if err != nil {
				log.Fatal(err)
			}

			// File & Folder Mode
			fMode := f.Mode()

			// Is folder
			if fMode.IsDir() {

				// Append to Folders Array
				folders = append(folders, path)

				// Is file
			} else if fMode.IsRegular() {

				// Append to Files Array
				files = append(files, path)
			}
		}

		return nil
	})

	return folders, files
}

func printStrings(slice []string) {
	f, err := os.OpenFile("testlogfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	for _, i := range slice {
		fmt.Fprintln(f, i)
	}

}
