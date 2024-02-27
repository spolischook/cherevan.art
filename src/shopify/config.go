package shopify

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/kpango/glg"
	"golang.org/x/oauth2"
)

type config struct {
	ShopName string `env:"SHOPIFY_SHOP_NAME,required"`
	ApiKey   string `env:"SHOPIFY_API_KEY,required"`
	ApiSecret string `env:"SHOPIFY_API_SECRET,required"`
	RedirectUrl string `env:"SHOPIFY_REDIRECT_URL,required"`
}

func Config() *config {
	if err := godotenv.Load(); err != nil {
		glg.Warn(err)
	}

	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg
}

func OauthConfig(cfg *config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ApiKey,
		ClientSecret: cfg.ApiSecret,
		RedirectURL:  cfg.RedirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + cfg.ShopName + "/admin/oauth/authorize",
			TokenURL: "https://" + cfg.ShopName + "/admin/oauth/access_token",
		},
	}
}
