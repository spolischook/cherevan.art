package shopify

import (
	goshopify "github.com/bold-commerce/go-shopify/v3"
)

func GetAllProducts(c *goshopify.Client) ([]goshopify.Product, error) {
	options := goshopify.ListOptions{Limit: 250}
	shopifyProducts := []goshopify.Product{}

	for {
		products, pagination, err := c.Product.ListWithPagination(options)
		if err != nil {
			return nil, err
		}
		shopifyProducts = append(shopifyProducts, products...)
		if pagination.NextPageOptions == nil {
			break
		}
		options.SinceID = products[len(products)-1].ID
	}

	return shopifyProducts, nil
}
