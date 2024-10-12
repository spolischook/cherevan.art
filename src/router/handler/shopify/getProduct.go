package shopify

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProduct(c *gin.Context) {
	if h.Client == nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	spId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   "invalid id",
			"message": err.Error(),
		})
		return
	}

	p, err := h.Client.Product.Get(spId, nil)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	r, _ := json.MarshalIndent(p, "", "  ")
	c.String(200, string(r))
}
