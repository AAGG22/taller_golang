package api

import (
	"sales-api/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func InitRoutes(router *gin.Engine) {
	localStorage := internal.NewLocalStorage()
	httpClient := resty.New()
	ventaService := internal.NewService(localStorage, httpClient)

	salesHandler := NewVentaHandler(ventaService)

	router.POST("/sales", salesHandler.HandleCreate)
	router.PATCH("/sales/:id", salesHandler.HandleUpdate)
	router.GET("/sales", salesHandler.HandleSearch)

}
