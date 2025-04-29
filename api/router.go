package api

import (
	"sales-api/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func InitRoutes(router *gin.Engine) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	localStorage := internal.NewLocalStorage()
	httpClient := resty.New()

	ventaService := internal.NewService(localStorage, httpClient, logger)
	salesHandler := NewVentaHandler(ventaService, logger)

	router.POST("/sales", salesHandler.HandleCreate)
	router.PATCH("/sales/:id", salesHandler.HandleUpdate)
	router.GET("/sales", salesHandler.HandleSearch)

}
