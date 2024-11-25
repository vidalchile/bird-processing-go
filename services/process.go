package services

import (
	"log"
	"sync"

	"project/apis"
	"project/models"
)

type FinalResult struct {
	ID          any
	Description any
	Images      []models.Image
}

var finalResults []FinalResult
var mu sync.Mutex // Para proteger el acceso a finalResults

// GetFinalResults devuelve una copia de los resultados finales
func GetFinalResults() []FinalResult {
	mu.Lock()
	defer mu.Unlock()

	// Retornar una copia para evitar problemas de concurrencia
	resultsCopy := make([]FinalResult, len(finalResults))
	copy(resultsCopy, finalResults)
	return resultsCopy
}

// ProcessRecord procesa un registro llamando a las APIs X2 y X3
func ProcessRecord(record models.Record, wg *sync.WaitGroup) {

	/*Siempre que lances una gorutina, coloca wg.Done() al principio de la función usando defer.
	Esto asegura que, sin importar cómo termine la tarea (exitosamente o con error),
	el contador del WaitGroup disminuirá correctament*/
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

	// Almacenar el resultado de forma segura
	mu.Lock()
	finalResults = append(finalResults, FinalResult{
		ID:          record.ID,
		Description: description.Description,
		Images:      images,
	})
	mu.Unlock()

	// Log de resultados
	log.Printf("Registro %d procesado exitosamente.\n", record.ID)

}
