package api

import (
	"fmt"
	"net/http"
	"sales-api/internal"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VentaHandler struct {
	ventaService *internal.Service
	logger       *zap.Logger
}

func NewVentaHandler(service *internal.Service, logger *zap.Logger) *VentaHandler {
	return &VentaHandler{
		ventaService: service,
		logger:       logger,
	}
}

func (ventaHandler *VentaHandler) HandleCreate(ginContext *gin.Context) {
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
		ventaHandler.logger.Error("error al crear venta", zap.Error(err))

		return
	}

	ventaHandler.logger.Info("venta created", zap.Any("venta", venta))
	ginContext.JSON(http.StatusCreated, venta)
}

func (ventaHandler *VentaHandler) HandleGet(ginContext *gin.Context) {
	ventaId := ginContext.Param("id")

	venta, err := ventaHandler.ventaService.GetVenta(ventaId)
	if err != nil {
		ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ventaHandler.logger.Info("venta encontrada", zap.Any("venta", venta))
	ginContext.JSON(http.StatusOK, venta)

}

func (ventaHandler *VentaHandler) HandleUpdate(ginContext *gin.Context) {
	id := ginContext.Param("id")

	var input struct {
		Status string `json:"status"`
	}

	if err := ginContext.ShouldBindJSON(&input); err != nil {
		ventaHandler.logger.Warn("json inválido al actualizar venta", zap.String("id", id), zap.Error(err))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Entrada invalida"})
		return
	}

	if input.Status != "approved" && input.Status != "rejected" {
		ventaHandler.logger.Warn("estado inválido al actualizar venta", zap.String("id", id), zap.String("status", input.Status))
		ginContext.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Filtro de estado inválido. Debe ser 'pending', 'approved' o 'rejected'"},
		)
		return
	}

	updatedVenta, err := ventaHandler.ventaService.UpdateStatus(id, input.Status)
	if err != nil {
		switch err {
		case internal.ErrVentaNotFound:
			ventaHandler.logger.Warn("Venta no encontrada", zap.String("id", id))
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
		case internal.ErrInvalidTransition:
			ventaHandler.logger.Warn("transición inválida", zap.String("id", id), zap.String("status", input.Status))
			ginContext.JSON(http.StatusConflict, gin.H{"error": "Transición de estado inválida"})
		default:
			ventaHandler.logger.Error("error interno al actualizar venta", zap.String("id", id), zap.Error(err))
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		}
		return
	}

	ventaHandler.logger.Info("venta actualizada correctamente",
		zap.String("id", updatedVenta.ID),
		zap.String("nuevo_estado", updatedVenta.Status))

	response := gin.H{
		"id":         updatedVenta.ID,
		"status":     updatedVenta.Status,
		"updated_at": updatedVenta.UpdatedAt.Format(time.RFC3339),
		"version":    updatedVenta.Version,
	}

	ginContext.JSON(http.StatusOK, response)
}

func (ventaHandler *VentaHandler) HandleSearch(ginContext *gin.Context) {
	userID := ginContext.Query("user_id")
	status := ginContext.Query("status")

	if userID == "" {
		ventaHandler.logger.Warn("falta user_id en búsqueda de ventas")
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	if status != "" && status != "pending" && status != "approved" && status != "rejected" {
		ventaHandler.logger.Warn("filtro de estado inválido en búsqueda", zap.String("status", status))
		ginContext.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Filtro de estado inválido. Debe ser 'pending', 'approved' o 'rejected'"},
		)
		return
	}

	ventas, err := ventaHandler.ventaService.SearchVentas(userID, status)
	if err != nil {
		ventaHandler.logger.Error("error al buscar ventas", zap.String("user_id", userID), zap.Error(err))
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	ventaHandler.logger.Info("búsqueda de ventas realizada con éxito",
		zap.String("user_id", userID),
		zap.String("status", status),
		zap.Int("cantidad", len(ventas)),
	)

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
