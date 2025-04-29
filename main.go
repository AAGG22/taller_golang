package main

import (
	"sales-api/api"
	"sales-api/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func main() {
	r := gin.Default()

	httpClient := resty.New()
	localStorage := internal.NewLocalStorage()
	api.InitRoutes(r, httpClient, localStorage)

	r.Run(":8081")

}
