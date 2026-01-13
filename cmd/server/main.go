package main

import (
	"log"

	"github.com/olsencastillo051172/forged-lro/api"
)

func main() {
	log.Println("FORGED-LRO server starting")

	// Registrar rutas de la API
	api.RegisterRoutes()

	// Aquí luego irá el arranque real del servidor
	// ej: http.ListenAndServe(":8080", nil)
}
