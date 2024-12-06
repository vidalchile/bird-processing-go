package main

import (
	"log"              // Paquete para registrar mensajes de log.
	"project/apis"     // Paquete para acceder a la API que obtiene los registros de aves.
	"project/models"   // Paquete donde se definen los modelos de datos (en este caso, 'Bird').
	"project/services" // Paquete donde se encuentran las funciones para procesar los registros.
	"sync"             // Paquete para usar WaitGroup y sincronización.
)

const maxConcurrentRequests = 5 // Número máximo de solicitudes concurrentes que pueden ejecutarse al mismo tiempo.

func main() {
	log.Println("Iniciando procesamiento...") // Mensaje de inicio del procesamiento.

	// Paso 1: Obtener registros de aves de la API "Chilean Birds".
	records, err := apis.GetBirds() // Llama a la API para obtener los registros de aves.
	if err != nil {                 // Si ocurre un error al obtener los registros, se registra el error y termina el programa.
		log.Fatalf("Error al obtener registros: %v\n", err)
	}

	var processedCount int // Contador para llevar la cuenta de cuántos registros han sido procesados.

	// Paso 2: Usar un grupo de trabajadores limitados para procesar los registros de forma concurrente.
	var wg sync.WaitGroup // Creamos una instancia de WaitGroup para esperar que todas las goroutines terminen.

	// Canal para controlar el número de solicitudes concurrentes.
	sem := make(chan struct{}, maxConcurrentRequests) // El canal 'sem' actúa como un semáforo. Limita la cantidad de goroutines concurrentes.

	// Iteramos sobre todos los registros de aves obtenidos desde la API.
	for _, record := range records {
		log.Println(">>> nombree: ", record.Name.Latin) // Registramos el nombre de la ave en formato latino.

		// Aseguramos que el canal tenga espacio para una nueva solicitud. Si no hay espacio, el programa espera.
		sem <- struct{}{} // Enviamos un valor vacío al canal. Esto actúa como un semáforo.

		wg.Add(1) // Indicamos que estamos esperando una nueva goroutine. Aumentamos el contador de WaitGroup.

		// Lanzamos una goroutine para procesar el registro de la ave de forma concurrente.
		go func(record models.Bird) {
			defer wg.Done() // Cuando esta goroutine termine, decrementamos el contador del WaitGroup.

			// Llamamos a la función 'ProcessRecord' para procesar el registro de la ave.
			services.ProcessRecord(record)

			// Liberamos el canal para permitir que otra goroutine comience.
			<-sem // Recibimos un valor del canal. Esto libera el semáforo y permite que otro proceso comience.
		}(record)

		processedCount++ // Aumentamos el contador de registros procesados.
	}

	// Esperamos a que todas las goroutines terminen su ejecución.
	wg.Wait() // El programa principal espera a que todas las goroutines en el WaitGroup terminen.

	log.Println("Procesamiento completado.") // Mensaje final después de que todas las goroutines hayan terminado.

	// Imprimimos el total de aves que fueron procesadas.
	log.Printf("Total de aves procesadas: %d", processedCount)
}
