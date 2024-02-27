package router

import (
	"github.com/cherevan.art/src/router/handler"
	"github.com/cherevan.art/src/router/handler/shopify"
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

	r.GET("/", handler.Root)
	r.GET("/shopify/login", shopifyHandler.Login)
	r.GET("/shopify/callback", shopifyHandler.Callback)
	r.POST("/shopify/update", shopifyHandler.UpdateProducts)

	return r
}
