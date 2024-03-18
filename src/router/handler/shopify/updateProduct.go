package shopify

import (
	"github.com/cherevan.art/src/artWork"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (h *Handler) UpdateProduct(c *gin.Context) {
	if h.Client == nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	spId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid id"})
		return
	}

	p, err := h.Client.Product.Get(spId, nil)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	aw, err := artWork.GetArtWorkByShopifyID(p.ID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	aw.UpdateShopifyProduct(p)

	p, err = h.Client.Product.Update(*p)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	aw.ShopifyOptionID = p.Variants[0].ID
	err = aw.Save()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "product updated"})
}
