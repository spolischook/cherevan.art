package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gopkg.in/yaml.v2"
)

type ArtWork struct {
	Title      string `yaml:"title"`
	Categories string
	InStock    bool
	IsVisible  bool
	Height     int
	Width      int
	Year       time.Time `yaml:"date"`
	Materials  []string
	Price      int
	ImageName  string
}

func main() {
	// Open the file
	csvfile, err := os.Open("cmd/import-art-works/inventory.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	r.Read() // skip header

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		artWork := parseRow(record)

		createArtWorkPage(artWork)
		fmt.Println("Created page for", artWork.Title)
	}
}

func createArtWorkPage(artWork ArtWork) {
	hugoContentDir := "content/art-works"
	if artWork.Title == "" {
		return
	}
	slug := Slugify(artWork.Title)
	contentDir := filepath.Join(hugoContentDir, slug)
	frontMatter := CreateFrontMatter(artWork)
	os.MkdirAll(contentDir, os.ModePerm)

	// Create Hugo Markdown file
	f, err := os.Create(filepath.Join(contentDir, "index.md"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Write the front matter to the file
	f.WriteString("---\n")
	f.WriteString(frontMatter)
	f.WriteString("---\n\n")

	// Write the content to the file
	//f.WriteString(fulltext)
}

func parseYear(year string) time.Time {
	if year == "" {
		return time.Date(1975, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	yearTime, err := time.Parse("2006", year)
	if err != nil {
		log.Printf(`Error parsing year: "%s"`, year)
		log.Println()
		log.Fatal(err)
	}
	return yearTime
}

func parseMaterials(materials string) []string {
	return strings.Split(materials, ",")
}

func parsePrice(price string) int {
	if price == "" {
		return -1
	}
	priceInt, err := strconv.Atoi(price)
	if err != nil {
		log.Printf(`Error parsing price: "%s"`, price)
		log.Println()
		log.Fatal(err)
	}
	return priceInt
}

func parseHeight(height string) int {
	if height == "" {
		return -1
	}
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		log.Printf(`Error parsing height: "%s"`, height)
		log.Println()
		log.Fatal(err)
	}
	return heightInt
}

func parseWidth(width string) int {
	if width == "" {
		return -1
	}
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		log.Printf(`Error parsing width: "%s"`, width)
		log.Println()
		log.Fatal(err)
	}
	return widthInt
}

func parseInStock(inStock string) bool {
	if inStock == "1" {
		return true
	}
	return false
}

func parseIsVisible(isVisible string) bool {
	if isVisible == "1" {
		return true
	}
	return false
}

func parseRow(row []string) ArtWork {
	title := row[2]
	cat := row[3]
	inStock := row[4]
	isVisible := row[5]
	height := row[6]
	width := row[7]
	year := row[8]
	materials := row[11]
	price := row[12]
	imageName := row[13]

	return ArtWork{
		Title:      title,
		Categories: cat,
		InStock:    parseInStock(inStock),
		IsVisible:  parseIsVisible(isVisible),
		Height:     parseHeight(height),
		Width:      parseWidth(width),
		Year:       parseYear(year),
		Materials:  parseMaterials(materials),
		Price:      parsePrice(price),
		ImageName:  imageName,
	}
}

func CreateFrontMatter(artwork ArtWork) string {
	// Convert the ArtWork struct to YAML
	yamlArtwork, err := yaml.Marshal(&artwork)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Convert the YAML to a string and return it
	return string(yamlArtwork)
}

func Slugify(s string) string {
	var slug string
	isAllowed := func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return r
		}
		if unicode.IsSpace(r) {
			return '-'
		}
		return -1
	}

	slug = strings.Map(isAllowed, strings.ToLower(s))

	// Remove any double dashes caused by spaces next to disallowed chars
	reg, _ := regexp.Compile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	return slug
}
