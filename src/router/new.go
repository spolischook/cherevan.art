package router

import (
	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/kpango/glg"

	"github.com/cherevan.art/src/router/handler"
	"github.com/cherevan.art/src/router/handler/shopify"
	sopifyLib "github.com/cherevan.art/src/shopify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(
		cors.Default(),
		gin.Logger(),
		gin.Recovery())
	r.LoadHTMLGlob("src/templates/**/*")

	shopifyHandler := shopify.Handler{}
	attachClient(&shopifyHandler)

	r.GET("/", handler.Root)
	r.GET("/shopify/login", shopifyHandler.Login)
	r.GET("/shopify/callback", shopifyHandler.Callback)
	r.POST("/shopify/update", shopifyHandler.UpdateProducts)
	r.POST("/shopify/update/:id", shopifyHandler.UpdateProduct)
	r.GET("/shopify/products/:id", shopifyHandler.GetProduct)

	return r
}

func attachClient(h *shopify.Handler) {
	config := sopifyLib.Config()

	if config.AccessToken == "" {
		glg.Warn("No access token for Shopify App, oAuth will be used.")
		return // should be logged in as a development app through oAuth
	}

	app := sopifyLib.NewApp(config)
	client := app.NewClient(
		config.ShopName,
		config.AccessToken,
		goshopify.WithLogger(&goshopify.LeveledLogger{}))

	h.Client = client
}
