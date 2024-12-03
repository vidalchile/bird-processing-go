package services

import (
	"log"
	"sync"

	"project/apis"
	"project/models"
)

type BirdName struct {
	Spanish string `json:"spanish"`
	English string `json:"english"`
	Latin   string `json:"latin"`
}

type FinalResultWikipedia struct {
	Name          BirdName `json:"name"`
	DataWikipedia any
}

var finalResults []FinalResultWikipedia
var mu sync.Mutex // Para proteger el acceso a finalResults

// GetFinalResults devuelve una copia de los resultados finales
func GetFinalResults() []FinalResultWikipedia {
	mu.Lock()
	defer mu.Unlock()

	// Retornar una copia para evitar problemas de concurrencia
	resultsCopy := make([]FinalResultWikipedia, len(finalResults))
	copy(resultsCopy, finalResults)
	return resultsCopy
}

// ProcessRecord procesa un registro llamando a las APIs X2 y X3
func ProcessRecord(record models.Bird, wg *sync.WaitGroup) {

	/*Siempre que lances una gorutina, coloca wg.Done() al principio de la función usando defer.
	Esto asegura que, sin importar cómo termine la tarea (exitosamente o con error),
	el contador del WaitGroup disminuirá correctament*/
	defer wg.Done()

	// Llamar a la API WIKIPEDIA
	wikipediaExtract, err := apis.CallWikipediaAPI(record.Name.Latin)
	if err != nil {
		log.Printf("Error en API WikipediaAP para registro %s: %v\n", record.UID, err)
		return
	}

	// Almacenar el resultado de forma segura
	mu.Lock()
	finalResults = append(finalResults, FinalResultWikipedia{
		Name:          BirdName(record.Name),
		DataWikipedia: wikipediaExtract,
	})
	mu.Unlock()

	// Log de resultados
	log.Printf("Registro %s procesado exitosamente.\n", record.UID)

}
