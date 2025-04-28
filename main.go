package main

import (
	"sales-api/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api.InitRoutes(r)

	r.Run(":8081")

}
