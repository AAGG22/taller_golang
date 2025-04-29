package api

import (
	"sales-api/internal"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func InitRoutes(router *gin.Engine, httpClient *resty.Client, localStorage *internal.LocalStorage) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	userService := internal.NewUserService(httpClient)
	ventaService := internal.NewService(localStorage, userService, logger)
	salesHandler := NewVentaHandler(ventaService, logger)

	router.POST("/sales", salesHandler.HandleCreate)
	router.PATCH("/sales/:id", salesHandler.HandleUpdate)
	router.GET("/sales", salesHandler.HandleSearch)
	router.GET("/sales/:id", salesHandler.HandleGet)
}
