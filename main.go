package main

import (
	"log"
	"sync"

	"project/apis"
	"project/services"
)

func main() {
	log.Println("Iniciando procesamiento...")

	// Paso 1: Obtener registros de la API X1
	records, err := apis.SimulateAPIX1()
	if err != nil {
		log.Fatalf("Error al obtener registros: %v\n", err)
	}

	// Paso 2: Procesar registros concurrentemente
	var wg sync.WaitGroup
	for _, record := range records {
		wg.Add(1)
		go services.ProcessRecord(record, &wg)
	}
	wg.Wait()

	log.Println("Procesamiento completado.")
}
