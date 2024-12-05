package services

import (
	"fmt"
	"log"
	"project/apis"
	"project/models"
)

// ProcessRecord procesa un registro llamando a las APIs X2 y X3
func ProcessRecord(record models.Bird) error {
	// Obtener detalle de la ave
	birdDetail, err := getBirdDetail(record)
	if err != nil {
		return fmt.Errorf("Error obteniendo detalles de la ave: %v", err)
	}

	// Llamar a la API WIKIPEDIA
	// birdExtract, err := getBirdExtract(record.Name.Latin)
	if err != nil {
		return fmt.Errorf("Error obteniendo datos de Wikipedia: %v", err)
	}

	// Log de resultados
	logResults(record, birdDetail, "")
	return nil
}

// getBirdDetail obtiene los detalles de una ave desde la API
func getBirdDetail(record models.Bird) (*apis.ResultBirdDetail, error) {
	detailBird, err := apis.GetBirdDetail(record.Links.Self)
	if err != nil {
		return nil, err
	}
	return &detailBird, nil
}

// getBirdExtract obtiene la información de Wikipedia de la ave
func getBirdExtract(latinName string) (string, error) {
	wikipediaExtract, err := apis.GetBirdExtract(latinName)
	if err != nil {
		return "", err
	}
	return wikipediaExtract, nil
}

// logResults maneja el logging de los resultados obtenidos
func logResults(record models.Bird, birdDetail *apis.ResultBirdDetail, birdExtract string) {
	log.Printf("Procesando registro %s...\n", record.UID)
	log.Printf("Nombre en español: %s, Nombre en inglés: %s, Nombre en latín: %s", record.Name.Spanish, record.Name.English, record.Name.Latin)
	log.Println("Mapa: ", birdDetail.Map)
	log.Println("IUCN: ", birdDetail.Iucn)
	log.Println("Migración: ", birdDetail.Migration)
	log.Println("Dimorfismo: ", birdDetail.Dimorphism)
	log.Println("Tamaño: ", birdDetail.Size)
	log.Println("Orden: ", birdDetail.Order)
	log.Println("Especie: ", birdDetail.Species)
	log.Println("Imágenes: ", birdDetail.Images)
	log.Println("Audio: ", birdDetail.Audio)
	log.Printf("Registro %s procesado exitosamente.\n", record.UID)
	log.Println("-------------------------------------------------------------------")
}
