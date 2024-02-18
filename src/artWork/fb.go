package artWork

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type fb struct {
	ID                       int    `csv:"id"`
	Title                    string `csv:"title"`
	Description              string `csv:"description"`
	Availability             string `csv:"availability"`
	Condition                string `csv:"condition"`
	Price                    string `csv:"price"`
	Link                     string `csv:"link"`
	ImageLink                string `csv:"image_link"`
	Brand                    string `csv:"brand"`
	GoogleProductCategory    string `csv:"google_product_category"`
	FbProductCategory        string `csv:"fb_product_category"`
	QuantityToSellOnFacebook string `csv:"quantity_to_sell_on_facebook"`
	SalePrice                string `csv:"sale_price"`
	SalePriceEffectiveDate   string `csv:"sale_price_effective_date"`
	ItemGroupId              string `csv:"item_group_id"`
	Gender                   string `csv:"gender"`
	Color                    string `csv:"color"`
	Size                     string `csv:"size"`
	AgeGroup                 string `csv:"age_group"`
	Material                 string `csv:"material"`
	Pattern                  string `csv:"pattern"`
	Shipping                 string `csv:"shipping"`
	ShippingWeight           string `csv:"shipping_weight"`
}

func (f fb) ToCsv() []string {
t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)
	row := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		row[i] = fmt.Sprintf("%v", field.Interface())
	}
	return row
}

func FbCsvHeaders() []string {
	t := reflect.TypeOf(fb{})
	headers := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		headers[i] = field.Tag.Get("csv")
	}
	return headers
}

func (a ArtWork) ToFb() fb {
	return fb{
		ID:                       a.ID,
		Title:                    a.Title,
		Description:              fbText(a),
		Availability:             fbAvailability(a),
		Condition:                "new",
		Price:                    fbPrice(a),
		Link:                     fmt.Sprintf("https://cherevan.art/%s/", a.Url),
		ImageLink:                fmt.Sprintf("https://cherevan.art/%s/%s", a.Url, a.ImageName),
		Brand:                    "CherevanArt",
		GoogleProductCategory:    "Home & Garden > Decor > Artwork",
		FbProductCategory:        "home > home goods > home decor",
		QuantityToSellOnFacebook: "1",
		SalePrice:                "",
		SalePriceEffectiveDate:   "",
		ItemGroupId:              strconv.Itoa(a.ID),
		Gender:                   "unisex",
		Color:                    "",
		Size:                     "",
		AgeGroup:                 "",
		Material:                 strings.Join(a.Materials, ", "),
		Pattern:                  "",
		Shipping:                 "",
		ShippingWeight:           "",
	}
}

func fbText(a ArtWork) string {
	if a.Text == "" {
		return "Original art work by Tetiana Cherevan"
	}

	if len(a.Text) > 9999 {
		return a.Text[:9999]
	}

	return a.Text
}

func fbAvailability(a ArtWork) string {
	if a.InStock {
		return "in stock"
	}
	return "out of stock"
}

func fbPrice(a ArtWork) string {
	return fmt.Sprintf("%s.00 EUR", a.Price)
}
