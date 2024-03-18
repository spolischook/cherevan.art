package shopify

import (
	"fmt"
	"net/http"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/gin-gonic/gin"

	"github.com/cherevan.art/src/shopify"
)

type Handler struct {
	Client *goshopify.Client
}

func (h *Handler) Login(c *gin.Context) {
	config := shopify.Config()
	app := shopify.NewApp(config)

	state := "nonce"
	authUrl := app.AuthorizeUrl(config.ShopName, state)
	c.Redirect(http.StatusFound, authUrl)
}

func (h *Handler) Callback(c *gin.Context) {
	config := shopify.Config()
	app := shopify.NewApp(config)

	if ok, _ := app.VerifyAuthorizationURL(c.Request.URL); !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid signature"))
		return
	}

	shopName := c.Query("shop")
	code := c.Query("code")
	token, err := app.GetAccessToken(shopName, code)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid token"))
		return
	}

	// Create a new API Client
	h.Client = app.NewClient(
		config.ShopName,
		token,
		goshopify.WithLogger(&goshopify.LeveledLogger{}))

	c.Redirect(http.StatusFound, "/")
}
