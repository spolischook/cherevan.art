package main

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/cherevan.art/src/shopify"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"net/http"
)

var app *goshopify.App

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config := shopify.Config()
	app = shopify.NewApp(config)

	state := "nonce"
	authUrl := app.AuthorizeUrl(config.ShopName, state)
	fmt.Println(authUrl)
	r := gin.Default()
	r.GET("/shopify/callback", func(c *gin.Context) {
		if ok, _ := app.VerifyAuthorizationURL(c.Request.URL); !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid signature"))
			return
		}

		shopName := c.Query("shop")
		code := c.Query("code")
		token, err := app.GetAccessToken(shopName, code)
		checkErr(err)

		// Create a new API client
		client := app.NewClient(
			config.ShopName,
			token,
			goshopify.WithLogger(&goshopify.LeveledLogger{}))

		{
			toks, err := client.StorefrontAccessToken.List(nil)
			checkErr(err)
			glg.Infof("There are %d tokens", len(toks))
			c.JSON(http.StatusOK, toks)
			return
		}
		{
			tok, err := client.StorefrontAccessToken.Create(goshopify.StorefrontAccessToken{
				Title: "Public Website",
			})
			checkErr(err)
			c.JSON(http.StatusOK, gin.H{
				"token": tok,
			})
			return
		}

		{
			// Fetch the number of products.
			numProducts, err := client.Product.Count(nil)
			checkErr(err)
			glg.Infof("There are %d products", numProducts)
			c.JSON(http.StatusOK, gin.H{
				"token":       token,
				"err":         err,
				"numProducts": numProducts,
			})
		}
	})
	r.Run(":8080")
}
