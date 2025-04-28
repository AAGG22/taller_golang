package main

import (
	"sales-api/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Inicializar las rutas
	api.InitRoutes(r)

	// Iniciar el servidor en el puerto 8080
	r.Run(":8081")

}
