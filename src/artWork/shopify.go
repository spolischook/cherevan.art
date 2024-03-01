package artWork

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/shopspring/decimal"
	"strings"
)

func (aw ArtWork) ShopifyProduct() goshopify.Product {
	return goshopify.Product{
		Title: aw.Title,
		Vendor: "CherevanArt",
		Tags: strings.Join(aw.Categories, ","),
		Status: "active",
		Variants: []goshopify.Variant{aw.ShopifyOrigVariant()},
	}
}

func (aw ArtWork) ShopifyOrigVariant() goshopify.Variant {
	price := decimal.NewFromInt(int64(aw.Price))
	return goshopify.Variant{
		Title: "Original",
		ProductID: aw.ShopifyID,
		Price: &price,
		InventoryPolicy: goshopify.VariantInventoryPolicyDeny,
		FulfillmentService: "manual",
		InventoryManagement: "shopify",
		InventoryQuantity: 1,
		RequireShipping: true,
	}
}

func (aw ArtWork) ShopifyImage() goshopify.Image {
	return goshopify.Image{
		ProductID: aw.ShopifyID,
		VariantIds: []int64{aw.ShopifyOptionID},
		Src: aw.ImageUrl(),
	}
}

func (aw ArtWork) UpdateShopifyProduct(sp *goshopify.Product) {
	sp.Title = aw.Title
	sp.Tags = strings.Join(aw.Categories, ", ")
	sp.Status = "active"
	size := fmt.Sprintf("%d x %d cm", aw.Width, aw.Height)
	medium := strings.Join(aw.Materials, ", ")
	sp.Options = []goshopify.ProductOption{
		{
			Name: "Type",
			Position: 1,
			Values: []string{"Original"},
		},
		{
			Name: "Size",
			Position: 2,
			Values: []string{size},
		},
		{
			Name: "Medium",
			Position: 3,
			Values: []string{medium},
		},
	}
	if len(sp.Variants) > 0 {
		sp.Variants[0] = aw.ShopifyOrigVariant()
	} else {
		sp.Variants = []goshopify.Variant{aw.ShopifyOrigVariant()}
	}
	sp.Variants[0].Option1 = "Original"
	sp.Variants[0].Option2 = size
	sp.Variants[0].Option3 = medium
	sp.Variants[0].ImageID = sp.Image.ID
}

func GetArtWorkByShopifyID(id int64) (*ArtWork, error) {
	aws, err := GetExistingArtWorks()
	if err != nil {
		return nil, err
	}

	for _, p := range aws {
		if p.ShopifyID == id {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("art work with shopify id %d not found", id)
}
