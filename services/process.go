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
	Map           any      `json:"map"`
	Iucn          any      `json:"iucn"`
	Migration     bool     `json:"migration"`
	Dimorphism    bool     `json:"dimorphism"`
	Size          string   `json:"size"`
	Order         string   `json:"order"`
	Species       string   `json:"species"`
	Images        any      `json:"images"`
	Audio         any      `json:"audio"`
	DataWikipedia any
}

var finalResults []FinalResultWikipedia
var mu sync.Mutex // Para proteger el acceso a finalResults

// ProcessRecord procesa un registro llamando a las APIs X2 y X3
func ProcessRecord(record models.Bird, wg *sync.WaitGroup) {

	/*Siempre que lances una gorutina, coloca wg.Done() al principio de la función usando defer.
	Esto asegura que, sin importar cómo termine la tarea (exitosamente o con error),
	el contador del WaitGroup disminuirá correctament*/
	defer wg.Done()

	// Obtener detalle de la ave
	detailBird, err := apis.CallDetailBirdAPI(record.Links.Self)
	if err != nil {
		log.Printf("Error en API Detail Bird para registro %s: %v\n", record.UID, err)
		return
	}

	// Llamar a la API WIKIPEDIA
	wikipediaExtract, err := apis.CallWikipediaAPI(record.Name.Latin)
	if err != nil {
		log.Printf("Error en API WikipediaAP para registro %s: %v\n", record.UID, err)
		return
	}

	log.Println("Names Bird: \n", BirdName(record.Name))
	log.Println("Map: \n", detailBird.Map)
	log.Println("Iucn: \n", detailBird.Iucn)
	log.Println("Migration: \n", detailBird.Migration)
	log.Println("Dimorphism: \n", detailBird.Dimorphism)
	log.Println("Size: \n", detailBird.Size)
	log.Println("Order: \n", detailBird.Order)
	log.Println("Species: \n", detailBird.Species)
	log.Println("Images: \n", detailBird.Images)
	log.Println("Audio: \n", detailBird.Audio)
	log.Println("Data Wikipedia: \n", wikipediaExtract)

	// Log de resultados
	log.Printf("Registro %s procesado exitosamente.\n", record.UID)

}
