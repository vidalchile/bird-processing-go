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
	var wg sync.WaitGroup // es una estructura de sincronización que permite esperar a que un grupo de gorutinas termine
	for _, record := range records {
		/*
			Antes de lanzar una gorutina, se incrementa el contador interno del WaitGroup en 1
			Esto indica que hay una nueva tarea (gorutina) pendiente
		*/
		wg.Add(1)

		//Se usa go para ejecutar la función services.ProcessRecord de forma concurrente
		//&wg: un puntero al WaitGroup para que cada gorutina pueda notificar cuando haya terminado su tarea
		go services.ProcessRecord(record, &wg)
	}
	wg.Wait() //Esperar a que todas las gorutinas finalicen

	// Imprimir los resultados finales
	log.Println("Resultados finales:", services.GetFinalResults())

	for _, result := range services.GetFinalResults() {
		log.Printf("ID: %d Descripción: %s Imágenes:\n", result.ID, result.Description)
		for _, img := range result.Images {
			log.Printf("  - %s\n", img)
		}
	}

	log.Println("Procesamiento completado.")
}
