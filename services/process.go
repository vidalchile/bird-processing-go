package services

import (
	"fmt"
	"log"
	"project/apis"
	"project/models"
)

// Result almacena los detalles del proceso de la ave y la información de Wikipedia.
type Result struct {
	BirdDetail  *apis.ResultBirdDetail // Detalles de la ave obtenidos de la API
	BirdExtract string                 // Extracto de Wikipedia para la ave
	Images      []apis.GetBirdImageDetailResult
	Error       error // Posible error durante el proceso
}

// ProcessRecord procesa un registro llamando a las APIs X2 y X3
// Esta función toma un registro de tipo 'Bird', llama a las APIs para obtener información
// adicional sobre la ave y luego simula almacenar estos datos en una base de datos.
func ProcessRecord(record models.Bird) error {
	var bufferCount = 3

	// Canal único para recibir los resultados: tanto el detalle de la ave como el extracto de Wikipedia.
	resultChannel := make(chan Result, bufferCount) // Canal con buffer de 2 para recibir ambos resultados

	// Goroutine para obtener los detalles de la ave.
	go func() {
		// Llamada a la API para obtener los detalles de la ave
		birdDetail, err := getBirdDetail(record)
		if err != nil {
			// Si ocurre un error al obtener los detalles, se envía el error al canal
			resultChannel <- Result{Error: fmt.Errorf("Error obteniendo detalles de la ave %s: %v", record.UID, err)}
			return
		}
		// Si se obtienen correctamente los detalles, se envían al canal
		resultChannel <- Result{BirdDetail: birdDetail}
	}()

	// Goroutine para obtener el extracto de Wikipedia relacionado con el nombre latino de la ave
	go func() {
		// Llamada a la API para obtener el extracto de Wikipedia de la ave
		birdExtract, err := getBirdExtract(record.Name.Latin)
		if err != nil {
			// Si ocurre un error al obtener el extracto de Wikipedia, se envía el error al canal
			log.Println("Error obteniendo extracto de Wikipedia para: ", record.Name.Latin)
			resultChannel <- Result{Error: fmt.Errorf("Error obteniendo extracto de Wikipedia para %s: %v", record.Name.Latin, err)}
		} else if birdExtract != "" {
			// Si se obtiene un extracto, se indica que la información de Wikipedia fue encontrada
			resultChannel <- Result{BirdExtract: fmt.Sprintf("Información de Wikipedia encontrada para %s", record.Name.Latin)}
		} else {
			log.Println("Error No se encontró información de Wikipedia para: ", record.Name.Latin)
			// Si no se encuentra información, se envía un mensaje indicando que no se halló
			resultChannel <- Result{BirdExtract: fmt.Sprintf("No se encontró información de Wikipedia para %s", record.Name.Latin)}
		}
	}()

	// Goroutine para obtener las imágenes asociadas
	go func() {
		// Títulos de imágenes asociados a la ave
		birdImages, err := apis.GetAllBirdImageDetails(record.Name.Latin)
		if err != nil {
			resultChannel <- Result{Error: fmt.Errorf("Error obteniendo imágenes de la ave %s: %v", record.Name.Latin, err)}
			return
		}
		// Agregamos las imágenes al resultado
		resultChannel <- Result{Images: birdImages}
	}()

	// Variables para almacenar los resultados recibidos del canal
	var birdDetail *apis.ResultBirdDetail
	var birdExtract string
	var birdImages []apis.GetBirdImageDetailResult
	var errorResult error

	// Recoger los resultados del canal
	// Se esperan dos resultados (detalles de la ave y extracto de Wikipedia)
	for i := 0; i < bufferCount; i++ {
		// Se recibe un resultado del canal
		result := <-resultChannel

		// Si el resultado contiene un error, se guarda para su posterior manejo
		if result.Error != nil {
			errorResult = result.Error
		}
		// Si se obtienen los detalles de la ave, se guardan en la variable correspondiente
		if result.BirdDetail != nil {
			birdDetail = result.BirdDetail
		}

		// Si se obtiene el extracto de Wikipedia, se guarda en la variable correspondiente
		if result.BirdExtract != "" {
			birdExtract = result.BirdExtract
		}

		// Si se obtiene imagenes de Wikimedia, se guarda en la variable correspondiente
		if result.Images != nil {
			birdImages = result.Images
		}
	}

	// Cerramos el canal después de recibir ambos resultados
	close(resultChannel)

	// Si ocurrió un error en alguna de las tareas (obtener detalles o extracto), se lo registramos
	if errorResult != nil {
		log.Println("Error:", errorResult)
		return errorResult
	}

	// Simulamos el almacenamiento en la base de datos con los datos obtenidos
	simulateDBSave(record, Result{BirdDetail: birdDetail, BirdExtract: birdExtract, Images: birdImages})

	// Retornamos nil si no hubo errores
	return nil
}

// getBirdDetail obtiene los detalles de una ave desde la API
// Esta función hace una llamada a la API externa (a través de apis.GetBirdDetail)
// para obtener los detalles de la ave a partir del enlace proporcionado en el registro.
func getBirdDetail(record models.Bird) (*apis.ResultBirdDetail, error) {
	// Llamada a la API para obtener los detalles de la ave
	detailBird, err := apis.GetBirdDetail(record.Links.Self)
	if err != nil {
		// Si ocurre un error en la llamada a la API, lo retornamos
		return nil, err
	}
	// Si no hubo errores, retornamos los detalles de la ave
	return &detailBird, nil
}

// getBirdExtract obtiene la información de Wikipedia de la ave
// Llama a la API externa para obtener el extracto de Wikipedia usando el nombre latino de la ave.
func getBirdExtract(latinName string) (string, error) {
	// Llamada a la API para obtener el extracto de Wikipedia de la ave
	wikipediaExtract, err := apis.GetBirdExtract(latinName)
	if err != nil {
		// Si ocurre un error en la llamada a la API, lo retornamos
		return "", err
	}
	// Si no hubo errores, retornamos el extracto de Wikipedia
	return wikipediaExtract, nil
}

// simulateDBSave simula el almacenamiento de los datos en una base de datos
// Esta función es solo una simulación para mostrar cómo se almacenarían los datos
// En un entorno real, aquí se realizaría una inserción real en una base de datos.
func simulateDBSave(record models.Bird, result Result) {
	// Simulación de inserción en la base de datos: Aquí, simplemente imprimimos los valores.
	fmt.Println("Simulando almacenamiento en la base de datos...")
	// Información del registro de la ave
	fmt.Printf("UID: %s\n", record.UID)
	fmt.Printf("Nombre en español: %s\n", record.Name.Spanish)
	fmt.Printf("Nombre en inglés: %s\n", record.Name.English)
	fmt.Printf("Nombre en latín: %s\n", record.Name.Latin)

	// Si los detalles de la ave están disponibles, los mostramos
	if result.BirdDetail != nil {
		fmt.Println("Detalles de la ave:")
		fmt.Printf("Mapa: %s\n", result.BirdDetail.Map)
		fmt.Printf("IUCN: %s\n", result.BirdDetail.Iucn)
		fmt.Printf("Migración: %v\n", result.BirdDetail.Migration)
		fmt.Printf("Dimorfismo: %v\n", result.BirdDetail.Dimorphism)
		fmt.Printf("Tamaño: %s\n", result.BirdDetail.Size)
		fmt.Printf("Orden: %s\n", result.BirdDetail.Order)
		fmt.Printf("Especie: %s\n", result.BirdDetail.Species)
		fmt.Println("Imágenes (API BIRDS): ", result.BirdDetail.Images)
		fmt.Println("Audio: ", result.BirdDetail.Audio)
	} else {
		// Si no se obtuvieron detalles, mostramos un mensaje indicando que no se encontraron
		fmt.Println("Detalles de la ave no encontrados.")
	}

	fmt.Println("Imágenes (Wikimedia): ", result.Images)

	// Mostramos el extracto de Wikipedia obtenido
	fmt.Println("Extracto de Wikipedia:", result.BirdExtract)
	fmt.Println("--------------------------------------------------")
}
