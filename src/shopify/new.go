package shopify

import goshopify "github.com/bold-commerce/go-shopify/v3"

func NewApp(c *config) *goshopify.App {
	return &goshopify.App{
		ApiKey:      c.ApiKey,
		ApiSecret:   c.ApiSecret,
		RedirectUrl: c.RedirectUrl,
		Scope:       Scopes(),
	}
}
