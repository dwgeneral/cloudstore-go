package main

import (
	"cloudstore-go/service/apigw/route"
)

func main() {
	r := route.Router()
	r.Run(":8080")
}
