package services

import (
	"log"
	"sync"

	"project/apis"
	"project/models"
)

// ProcessRecord procesa un registro llamando a las APIs X2 y X3
func ProcessRecord(record models.Record, wg *sync.WaitGroup) {
	defer wg.Done()

	// Paso 1: Llamar a la API X2
	description, err := apis.SimulateAPIX2(record.ID)
	if err != nil {
		log.Printf("Error en API X2 para registro %d: %v\n", record.ID, err)
		return
	}

	// Paso 2: Llamar a la API X3
	images, err := apis.SimulateAPIX3(record.ID)
	if err != nil {
		log.Printf("Error en API X3 para registro %d: %v\n", record.ID, err)
		return
	}

	// Log de resultados
	log.Printf("Registro %d procesado: %s\n", record.ID, description.Description)
	for _, img := range images {
		log.Printf("  Imagen: %s\n", img.URL)
	}
}
