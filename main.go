package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/signintech/gopdf"
)

func main() {
	var inputPath string
	flag.StringVar(&inputPath, "i", ".", "Input path to folder")
	var outputPath string
	flag.StringVar(&outputPath, "o", "output.pdf", "Output path for pdf")
	flag.Parse()

	folders := getSortedFolders(inputPath)

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for _, folder := range folders {
		files := getSortedFiles(inputPath + "/" + folder)
		for _, file := range files {
			fullpath := inputPath + "/" + folder + "/" + file
			fmt.Println(fullpath)

			pdf.AddPage()

			imageW, imageH := getImageDimension(fullpath)
			scale := math.Min(gopdf.PageSizeA4.W/float64(imageW), gopdf.PageSizeA4.H/float64(imageH))
			var rect gopdf.Rect
			rect.W = float64(imageW) * scale
			rect.H = float64(imageH) * scale

			centerW := (gopdf.PageSizeA4.W - rect.W) / 2
			centerH := (gopdf.PageSizeA4.H - rect.H) / 2

			err := pdf.Image(fullpath, centerW, centerH, &rect)
			if err != nil {
				log.Panic(err)
			}
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

type sortablePath struct {
	original  string
	numerical int
}

func makeSortable(entries []os.FileInfo, onlyFolders bool) []sortablePath {
	// Prepare buffer for the numerically named folders
	var buf []sortablePath

	for _, entry := range entries {
		if (entry.IsDir() && onlyFolders) || (!entry.IsDir() && !onlyFolders) {
			// Remove everything that isn't number, so "page420" reads as "420"
			reg, err := regexp.Compile("[^0-9]+")
			if err != nil {
				fmt.Println("Skipping {}", entry.Name())
				continue
			}

			onlyNumbers := reg.ReplaceAllString(entry.Name(), "")
			numeric, err := strconv.Atoi(onlyNumbers)
			if err != nil {
				fmt.Println("Skipping {}", entry.Name())
				continue
			}
			if len(onlyNumbers) > 0 {
				buf = append(buf, sortablePath{entry.Name(), numeric})
			}
		}
	}
	return buf
}

func sortPaths(s []sortablePath) []string {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].numerical < s[j].numerical
	})
	var onlyPaths []string
	for _, sp := range s {
		onlyPaths = append(onlyPaths, sp.original)
	}
	return onlyPaths
}

func getSortedFolders(dirPath string) []string {
	// Get all things in directory
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic("could not open folder " + err.Error())
	}

	sortable := makeSortable(entries, true)
	return sortPaths(sortable)
}

func getSortedFiles(folder string) []string {
	entries, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	sortable := makeSortable(entries, false)
	return sortPaths(sortable)
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
