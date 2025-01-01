package main

import "github.com/cherevan.art/src/router"

func main() {
	r := router.New()
	r.Run(":9085")
}
