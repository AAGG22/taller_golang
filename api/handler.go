package api

import (
	"fmt"
	"net/http"
	"sales-api/internal"
	"time"

	"github.com/gin-gonic/gin"
)

type ventaHandler struct {
	ventaService *internal.Service
}

func NewVentaHandler(service *internal.Service) *ventaHandler {
	return &ventaHandler{
		ventaService: service,
	}
}

func (ventaHandler *ventaHandler) HandleCreate(ginContext *gin.Context) {
	var venta internal.Venta
	if err := ginContext.ShouldBindJSON(&venta); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Entrada invalida"})
		return
	}

	fmt.Printf("venta:%v\n", venta)

	err := ventaHandler.ventaService.Create(&venta)
	if err != nil {
		if err == internal.ErrInvalidAmount || err == internal.ErrUserNotFound {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, venta)
}

func (ventaHandler *ventaHandler) HandleUpdate(ginContext *gin.Context) {
	id := ginContext.Param("id")

	var input struct {
		Status string `json:"status"`
	}

	if err := ginContext.ShouldBindJSON(&input); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Entrada invalida"})
		return
	}

	if input.Status != "approved" && input.Status != "rejected" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Filtro de estado inválido. Debe ser 'pending', 'approved' o 'rejected'"})
		return
	}

	updatedVenta, err := ventaHandler.ventaService.UpdateStatus(id, input.Status)
	if err != nil {
		switch err {
		case internal.ErrVentaNotFound:
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
		case internal.ErrInvalidTransition:
			ginContext.JSON(http.StatusConflict, gin.H{"error": "Transición de estado inválida"})
		default:
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		}
		return
	}

	response := gin.H{
		"id":         updatedVenta.ID,
		"status":     updatedVenta.Status,
		"updated_at": updatedVenta.UpdatedAt.Format(time.RFC3339),
		"version":    updatedVenta.Version,
	}

	ginContext.JSON(http.StatusOK, response)
}

func (ventaHandler *ventaHandler) HandleSearch(ginContext *gin.Context) {
	userID := ginContext.Query("user_id")
	status := ginContext.Query("status")

	if userID == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "user_id i es requerido"})
		return
	}

	if status != "" && status != "pending" && status != "approved" && status != "rejected" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Filtro de estado inválido. Debe ser 'pending', 'approved' o 'rejected'"})
		return
	}

	ventas, err := ventaHandler.ventaService.SearchVentas(userID, status)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	var quantity, approved, rejected, pending int
	var totalAmount float64

	for _, v := range ventas {
		quantity++
		totalAmount += v.Amount
		switch v.Status {
		case "approved":
			approved++
		case "rejected":
			rejected++
		case "pending":
			pending++
		}
	}

	response := gin.H{
		"metadata": gin.H{
			"quantity":     quantity,
			"approved":     approved,
			"rejected":     rejected,
			"pending":      pending,
			"total_amount": totalAmount,
		},
		"results": ventas,
	}

	ginContext.JSON(http.StatusOK, response)
}
