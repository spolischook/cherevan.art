package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/kpango/glg"
	"github.com/nfnt/resize"
	"gopkg.in/yaml.v2"

	"github.com/cherevan.art/src/googleDrive"
)

type ArtWork struct {
	Title      string    `yaml:"title" json:"title"`
	Slug       string    `yaml:"slug" json:"slug"`
	Categories []string  `yaml:"categories" json:"categories"`
	InStock    bool      `yaml:"inStock" json:"inStock"`
	IsVisible  bool      `yaml:"isVisible" json:"isVisible"`
	Height     int       `yaml:"height" json:"height"`
	Width      int       `yaml:"width" json:"width"`
	Date       time.Time `yaml:"date" json:"date"`
	Materials  []string  `yaml:"materials" json:"materials"`
	Price      int       `yaml:"price" json:"price"`
	ImageName  string    `yaml:"mainImage" json:"mainImage"`
}

const inventoryFileID = "1MPkRsuhJhJdYlTCOxHFiJoFX3DZalZPbTOUObacJPg0"
const paintingsFolderId = "1NA3FqpicdYnl7RS7-YDks0oh6eLNFTcp"
const graphicsFolderId = "11nADUmGwSlnT4LGA9DF7RzButDb1FHNd"

func checkErr(err error) {
	if err != nil {
		glg.Fatal(err)
	}
}

func main() {
	srv, err := googleDrive.CreateService()
	checkErr(err)
	err = srv.DownloadSheetAsCSV(inventoryFileID, "inventory.csv")
	checkErr(err)
	_ = srv

	// Open the file
	csvfile, err := os.Open("inventory.csv")
	checkErr(err)

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
		checkErr(err)

		artWork := parseRow(record)

		updateArtWorkPage(artWork)
		fetchArtWorkImage(srv, artWork)
		glg.Info("Created page for", artWork.Title)
	}
}

func scaleImage(src string, width uint) {
	// Open the image file
	file, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode the image
	img, err := jpeg.Decode(file)
	if err != nil {
		glg.Fatalf("Error decoding image %s: %s\n", src, err)
	}

	// Scale the image
	scaledImg := resize.Resize(width, 0, img, resize.Lanczos3)

	// Create the output file
	out, err := os.Create(src)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Write the scaled image to the output file
	jpeg.Encode(out, scaledImg, nil)
}

func findFileInDir(dirPath string, fileName string) (string, error) {
	var filePath string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == fileName {
			filePath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if filePath == "" {
		return "", fmt.Errorf("file not found: %s\n", fileName)
	}
	return filePath, nil
}

func fetchArtWorkImage(srv *googleDrive.Service, artWork ArtWork) {
	if artWork.ImageName == "" {
		glg.Warnf("No image name for art work: %#v\n", artWork)
		return
	}
	imageName := artWork.ImageName
	imageDir := getArtWorkDir(artWork.Title)

	// check if the image already exists
	_, err := findFileInDir(imageDir, imageName)
	if err == nil {
		return
	}

	// if artWork.Categories contains "painting" then use paintingsFolderId
	// else use graphicsFolderId
	var dirId string
	for _, cat := range artWork.Categories {
		if cat == "painting" {
			dirId = paintingsFolderId
			break
		}
		if cat == "graphics" {
			dirId = graphicsFolderId
			break
		}
	}

	MkdirAll(imageDir, os.ModePerm)
	err = srv.DownloadFile(imageName, dirId, imageDir)
	if err == nil {
		imagePath := filepath.Join(imageDir, imageName)
		scaleImage(imagePath, 1200)
		glg.Info("Downloaded image for", artWork.Title)
	} else {
		glg.Fatalf("Error downloading image %s: %s\n", imageName, err)
	}
}

func checkDestDir(dst string) error {
	const protectedDir = "/home/spolischook/cherryDrive"
	if strings.HasPrefix(dst, protectedDir) {
		err := fmt.Errorf("Cannot write to directory: %s\n", protectedDir)
		glg.Error(err)
		return err
	}

	return nil
}

func MkdirAll(path string, perm os.FileMode) error {
	if err := checkDestDir(path); err != nil {
		glg.Error(err)
		return err
	}

	return os.MkdirAll(path, perm)
}

func copyFile(dst, src string) error {
	if err := checkDestDir(dst); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		glg.Error(err)
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		glg.Error(err)
		return err
	}
	defer dstFile.Close()
	// Copy the content from the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		glg.Error(err)
		return err
	}

	return nil
}

func createFile(path string) (*os.File, error) {
	if err := checkDestDir(path); err != nil {
		return nil, err
	}

	return os.Create(path)
}

func updateArtWorkPage(artWork ArtWork) {
	if artWork.Title == "" {
		glg.Errorf("No title for art work: %#v\n", artWork)
		return
	}

	contentDir := getArtWorkDir(artWork.Title)
	MkdirAll(contentDir, os.ModePerm)

	// Open the existing Hugo Markdown file
	filePath := filepath.Join(contentDir, "index.md")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		glg.Fatal(err)
	}
	defer f.Close()

	// Read the file content
	content, err := io.ReadAll(f)
	if err != nil {
		glg.Fatal(err)
	}

	// Parse the existing front matter
	existingData := map[string]interface{}{}
	yaml.Unmarshal(content, &existingData)

	// Convert the ArtWork struct to a map
	newData, _ := structToMap(artWork)

	// Merge the existing front matter with the new data
	for key, value := range newData {
		existingData[key] = value
	}

	// Create the updated front matter
	frontMatter, _ := yaml.Marshal(existingData)

	// Write the updated front matter back to the file
	f, err = os.Create(filePath)
	if err != nil {
		glg.Fatal(err)
	}
	defer f.Close()

	f.WriteString("---\n")
	f.WriteString(string(frontMatter))
	f.WriteString("---\n\n")

	// Write the content to the file
	//f.WriteString(fulltext)
}

func structToMap(item interface{}) (map[string]interface{}, error) {
	b, _ := json.Marshal(item)
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m, nil
}

func createArtWorkPage(artWork ArtWork) {
	if artWork.Title == "" {
		glg.Errorf("No title for art work: %#v\n", artWork)
		return
	}

	contentDir := getArtWorkDir(artWork.Title)
	MkdirAll(contentDir, os.ModePerm)
	frontMatter := CreateFrontMatter(artWork)

	// Create Hugo Markdown file
	f, err := createFile(filepath.Join(contentDir, "index.md"))
	if err != nil {
		glg.Fatal(err)
	}
	defer f.Close()

	// Write the front matter to the file
	f.WriteString("---\n")
	f.WriteString(frontMatter)
	f.WriteString("---\n\n")

	// Write the content to the file
	//f.WriteString(fulltext)
}

func getArtWorkDir(title string) string {
	slug := Slugify(title)
	hugoContentDir := "content/art-works"
	contentDir := filepath.Join(hugoContentDir, slug)
	return contentDir
}

var timeCounters = map[string]time.Time{}

func parseYear(year string) time.Time {
	if year == "" {
		glg.Fatal("Year is empty")
	}
	if t, ok := timeCounters[year]; ok {
		t = t.Add(-time.Second)
		timeCounters[year] = t
		return t
	}
	yearTime, err := time.Parse("2006", year)
	if err != nil {
		glg.Printf(`Error parsing year: "%s"`, year)
		glg.Println()
		glg.Fatal(err)
	}
	t := time.Date(yearTime.Year(), time.December, 31, 23, 59, 59, 0, time.UTC)
	timeCounters[year] = t

	return t
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
		glg.Printf(`Error parsing price: "%s"`, price)
		glg.Println()
		glg.Fatal(err)
	}
	return priceInt
}

func parseHeight(height string) int {
	if height == "" {
		return -1
	}
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		glg.Printf(`Error parsing height: "%s"`, height)
		glg.Println()
		glg.Fatal(err)
	}
	return heightInt
}

func parseWidth(width string) int {
	if width == "" {
		return -1
	}
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		glg.Printf(`Error parsing width: "%s"`, width)
		glg.Println()
		glg.Fatal(err)
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

func parseCategories(categories string) []string {
	return strings.Split(categories, "/")
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
		Slug:       Slugify(title),
		Categories: parseCategories(cat),
		InStock:    parseInStock(inStock),
		IsVisible:  parseIsVisible(isVisible),
		Height:     parseHeight(height),
		Width:      parseWidth(width),
		Date:       parseYear(year),
		Materials:  parseMaterials(materials),
		Price:      parsePrice(price),
		ImageName:  imageName,
	}
}

func CreateFrontMatter(artwork ArtWork) string {
	// Convert the ArtWork struct to YAML
	yamlArtwork, err := yaml.Marshal(&artwork)
	if err != nil {
		glg.Fatalf("error: %v", err)
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
