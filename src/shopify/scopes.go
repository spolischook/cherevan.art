package shopify

import (
	"strings"
)

var scopes = []string{
	"read_products",
	"write_products",
	"read_orders",
	"write_orders",
	"read_inventory",
	"write_inventory",
	"unauthenticated_read_checkouts",
	"unauthenticated_write_checkouts",
	"unauthenticated_read_customers",
	"unauthenticated_write_customers",
	"unauthenticated_read_customer_tags",
	"unauthenticated_read_content",
	"unauthenticated_read_metaobjects",
	"unauthenticated_read_product_inventory",
	"unauthenticated_read_product_listings",
	"unauthenticated_read_product_pickup_locations",
	"unauthenticated_read_product_tags",
	"unauthenticated_read_selling_plans",
}

func Scopes() string {
	return strings.Join(scopes, ",")
}
