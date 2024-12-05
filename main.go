package main

import (
	"log"
	"sync"

	"project/apis"
	"project/models"
	"project/services"
)

const maxConcurrentRequests = 5 // Máximo número de solicitudes concurrentes

func main() {
	log.Println("Iniciando procesamiento...")

	// Paso 1: Obtener registros aves de la API "Chilean Birds"
	records, err := apis.GetBirds()
	if err != nil {
		log.Fatalf("Error al obtener registros: %v\n", err)
	}

	var processedCount int // Contador de procesos completados

	// Paso 2: Usar un grupo de trabajadores limitados para procesar los registros concurrentemente
	var wg sync.WaitGroup
	// Canal para limitar el número de solicitudes concurrentes
	sem := make(chan struct{}, maxConcurrentRequests)

	for _, record := range records {
		// Asegurarse de que el canal tiene espacio para una nueva solicitud
		sem <- struct{}{}
		wg.Add(1)

		// Lanzar la goroutine para procesar el registro
		go func(record models.Bird) {
			defer wg.Done()
			// Llamar a la función de procesamiento
			services.ProcessRecord(record)

			// Liberar el canal cuando termine
			<-sem
		}(record)

		processedCount++
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	log.Println("Procesamiento completado.")
	log.Printf("Total de aves procesadas: %d", processedCount)
}
