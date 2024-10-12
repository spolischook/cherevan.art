package shopify

import (
	"fmt"
	"github.com/cherevan.art/src/artWork"
	"github.com/cherevan.art/src/shopify"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *Handler) UpdateProducts(c *gin.Context) {
	if h.Client == nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	existingArtWorks, err := artWork.GetExistingArtWorks()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get all products from Shopify
	shopifyProducts, err := shopify.GetAllProducts(h.Client)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Update products
	for _, aw := range existingArtWorks {
		if aw.Salable() != nil {
			continue
		}

		// Check if product exists in Shopify
		isNew := true
		for _, sp := range shopifyProducts {
			if aw.ShopifyID == sp.ID {
				isNew = false
				aw.UpdateShopifyProduct(&sp)
				time.Sleep(time.Second / 2) // prevent exceeding the limit of requests
				sp1, err := h.Client.Product.Update(sp)
				if err != nil {
					c.AbortWithStatusJSON(
						500,
						gin.H{
							fmt.Sprintf(`error update %d %s:`, aw.ID, aw.Title): err.Error()},
					)
					return
				}
				aw.ShopifyOptionID = sp1.Variants[0].ID
				err = aw.Save()
				if err != nil {
					c.AbortWithStatusJSON(
						500,
						gin.H{
							fmt.Sprintf(`error save %d %s:`, aw.ID, aw.Title): err.Error()},
					)
					return
				}
			}
		}

		if !isNew {
			continue
		}

		// Create product
		createdProduct, err := h.Client.Product.Create(aw.ShopifyProduct())
		if err != nil {
			c.AbortWithStatusJSON(
				500,
				gin.H{
					fmt.Sprintf(`error create %d %s:`, aw.ID, aw.Title): err.Error()},
			)
			return
		}

		_, err = h.Client.ProductListing.Publish(createdProduct.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error publish %d %s: %s", aw.ID, aw.Title, err.Error())})
			return
		}

		aw.ShopifyID = createdProduct.ID
		aw.ShopifyOptionID = createdProduct.Variants[0].ID
		err = aw.Save()
		if err != nil {
			panic(err)
		}

		// Create image
		_, err = h.Client.Image.Create(createdProduct.ID, aw.ShopifyImage())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		}

		aw.UpdateShopifyProduct(createdProduct)
		createdProduct, err = h.Client.Product.Update(*createdProduct)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		}
	}

	// todo: we should delete only original artworks, but not all other products, e.g. t-shorts, prints, etc.
	//for _, p := range shopifyProducts {
	//	isDeleted := true
	//	for _, aw := range existingArtWorks {
	//		if p.ID == aw.ShopifyID {
	//			isDeleted = false
	//		}
	//	}
	//	if false == isDeleted {
	//		continue
	//	}
	//
	//	glg.Warnf("deleting product %d: %s\n", p.ID, p.Title)
	//	err := h.Client.Product.Delete(p.ID)
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	//	}
	//}
}
