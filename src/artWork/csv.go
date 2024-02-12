package artWork

import (
	"fmt"
	"github.com/cherevan.art/src/tool"
	"github.com/kpango/glg"
	"strconv"
	"strings"
	"time"
)

type csv struct {
	ID        int
	Title     int
	Category  int
	InStock   int
	IsVisible int
	Height    int
	Width     int
	Year      int
	Materials int
	Price     int
	ImageName int
}
type csvHeaders []string

func NewCsv(headers csvHeaders) *csv {
	return &csv{
		ID:        headers.MustGetIndex("ID"),
		Title:     headers.MustGetIndex("Title en"),
		Category:  headers.MustGetIndex("Category"),
		InStock:   headers.MustGetIndex("in stock"),
		IsVisible: headers.MustGetIndex("is visible"),
		Height:    headers.MustGetIndex("H"),
		Width:     headers.MustGetIndex("W"),
		Year:      headers.MustGetIndex("Year"),
		Materials: headers.MustGetIndex("Materials en"),
		Price:     headers.MustGetIndex("Price USD"),
		ImageName: headers.MustGetIndex("Image name"),
	}
}

func (c csv) ParseRow(row []string) ArtWork {
	id := row[c.ID]
	title := row[c.Title]
	cat := row[c.Category]
	inStock := row[c.InStock]
	isVisible := row[c.IsVisible]
	height := row[c.Height]
	width := row[c.Width]
	year := row[c.Year]
	materials := row[c.Materials]
	price := row[c.Price]
	imageName := row[c.ImageName]

	date, _ := parseYear(year)

	return ArtWork{
		ID:         parseInt(id),
		Title:      title,
		Slug:       tool.Slugify(title),
		Categories: parseCategories(cat),
		InStock:    parseBool(inStock),
		IsVisible:  parseBool(isVisible),
		Height:     parseInt(height),
		Width:      parseInt(width),
		Date:       date,
		Materials:  parseMaterials(materials),
		Price:      parseInt(price),
		ImageName:  imageName,
	}
}

func (hs csvHeaders) MustGetIndex(header string) int {
	for i, h := range hs {
		if h == header {
			return i
		}
	}

	// this will stop execution but compiler doesn't know that
	glg.Fatalf(`Header "%s" not found`, header)
	return -1
}

func parseYear(year string) (time.Time, error) {
	if year == "" {
		err := fmt.Errorf("year is empty")
		return time.Time{}, err
	}

	yearTime, err := time.Parse("2006", year)
	if err != nil {
		err := fmt.Errorf(`Error parsing year: "%s"`, year)
		glg.Error(err)
		return time.Time{}, err
	}

	return yearTime, nil
}

func parseMaterials(materials string) []string {
	return strings.Split(materials, ",")
}

func parseInt(n string) int {
	if n == "" {
		return -1
	}
	heightInt, err := strconv.Atoi(n)
	if err != nil {
		glg.Errorf(`Error parsing number: "%s"\n`, n)
	}
	return heightInt
}

func parseBool(v string) bool {
	if v == "1" {
		return true
	}
	return false
}

func parseCategories(categories string) []string {
	return strings.Split(categories, "/")
}
