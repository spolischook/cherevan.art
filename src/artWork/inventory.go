package artWork

import (
	"fmt"
	"github.com/cherevan.art/src/tool"
	"github.com/kpango/glg"
	"strconv"
	"strings"
	"time"
)

type inventory struct {
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
type inventoryHeaders []string

func NewCsvInventory(headers inventoryHeaders) *inventory {
	return &inventory{
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

func (inv inventory) ParseRow(row []string) ArtWork {
	id := row[inv.ID]
	title := row[inv.Title]
	cat := row[inv.Category]
	inStock := row[inv.InStock]
	isVisible := row[inv.IsVisible]
	height := row[inv.Height]
	width := row[inv.Width]
	year := row[inv.Year]
	materials := row[inv.Materials]
	price := row[inv.Price]
	imageName := row[inv.ImageName]

	date, _ := parseYear(year)

	aw := ArtWork{
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
	aw.HugoUrl = aw.GetUrl()
	return aw
}

func (hs inventoryHeaders) MustGetIndex(header string) int {
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
