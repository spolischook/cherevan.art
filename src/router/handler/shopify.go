package handler

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/cherevan.art/src/shopify"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"net/http"
)

func Callback(c *gin.Context) {
	config := shopify.Config()
	app := shopify.NewApp(config)

	if ok, _ := app.VerifyAuthorizationURL(c.Request.URL); !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid signature"))
		return
	}

	shopName := c.Query("shop")
	code := c.Query("code")
	token, err := app.GetAccessToken(shopName, code)

	// Create a new API client
	client := app.NewClient(
		config.ShopName,
		token,
		goshopify.WithLogger(&goshopify.LeveledLogger{}))

	// Fetch the number of products.
	numProducts, err := client.Product.Count(nil)

	if err != nil {
		glg.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the number of products: " + err.Error()})
	}

	glg.Infof("There are %d products", numProducts)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"err":   err,
		"numProducts": numProducts,
	})
}
