package shopify

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"slices"
	"strings"

	"github.com/cherevan.art/src/artWork"
)

func (h *Handler) UpdateProducts(c *gin.Context) {
	if h.client == nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	existingArtWorks, err := artWork.GetExistingArtWorks()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get all products from Shopify
	options := goshopify.ListOptions{Limit: 250}
	shopifyProducts := []goshopify.Product{}

	count, err := h.client.Product.Count(nil)
	_ = count

	for {
		products, pagination, err := h.client.Product.ListWithPagination(options)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
		shopifyProducts = append(shopifyProducts, products...)
		if pagination.NextPageOptions == nil {
			break
		}
		options.SinceID = products[len(products)-1].ID
	}

	shopifyProducts2 := slices.CompactFunc(shopifyProducts, func(a, b goshopify.Product) bool {
		return a.ID == b.ID
	})
	fmt.Println("Len of shopifyProducts and shopifyProducts2:", len(shopifyProducts), len(shopifyProducts2))

	// Update products
	for _, aw := range existingArtWorks {
		if aw.Salable() != nil {
			continue
		}

		// Check if product exists in Shopify
		isNew := true
		for _, sp := range shopifyProducts {
			if aw.ShopifyID == sp.ID {
				// Update product
				// ...
				isNew = false

			}
		}

		if false == isNew {
			continue
		}

		price := decimal.NewFromInt(int64(aw.Price))
		// Create product
		newProduct := goshopify.Product{
			Title: aw.Title,
			Vendor: "CherevanArt",
			Tags: strings.Join(aw.Categories, ", "),
			Status: "active",
			Variants: []goshopify.Variant{
				{
					ProductID: int64(aw.ID),
					Title: "Original Art Work",
					Price: &price,
					FulfillmentService: "manual",
					InventoryManagement: "shopify",
					InventoryPolicy: "deny",
					InventoryQuantity: 1,
					RequireShipping: true,
				},
			},
		}
		createdProduct, err := h.client.Product.Create(newProduct)
		if err != nil {
			c.AbortWithStatusJSON(
				500,
				gin.H{
					fmt.Sprintf(`error create %d %s:`, aw.ID, aw.Title): err.Error()},
				)
			return
		}

		aw.ShopifyID = createdProduct.ID
		aw.ShopifyOptionID = createdProduct.Variants[0].ID
		err = aw.Save()
		if err != nil {
			panic(err)
		}

		// Create image
		img := &goshopify.Image{
			ProductID: createdProduct.ID,
			VariantIds: []int64{createdProduct.Variants[0].ID},
			Src: aw.ImageUrl(),
		}

		_, err = h.client.Image.Create(createdProduct.ID, *img)
		if err != nil {
				c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			}
	}
}
