package artWork

import (
	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/shopspring/decimal"
	"strings"
)

func (aw ArtWork) ShopifyProduct() goshopify.Product {
	return goshopify.Product{
		Title: aw.Title,
		Vendor: "CherevanArt",
		Tags: strings.Join(aw.Categories, ", "),
		Status: "active",
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
