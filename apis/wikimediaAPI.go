package apis

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)

// --------------- LISTADO DE IMAGENES ---------------------------------------

// SearchResult representa cada resultado individual en la sección "search".
type SearchResult struct {
	NS        int       `json:"ns"`
	Title     string    `json:"title"`
	PageID    int       `json:"pageid"`
	Size      int       `json:"size"`
	WordCount int       `json:"wordcount"`
	Snippet   string    `json:"snippet"`
	Timestamp time.Time `json:"timestamp"`
}

// Query contiene la información de la consulta y la lista de resultados.
type Query struct {
	SearchInfo struct {
		TotalHits int `json:"totalhits"`
	} `json:"searchinfo"`
	Search []SearchResult `json:"search"`
}

// Continue contiene los datos para la paginación.
type Continue struct {
	SrOffset int    `json:"sroffset"`
	Continue string `json:"continue"`
}

// APIResponse representa la estructura completa del JSON de respuesta.
type WikimediaAPIResponse struct {
	BatchComplete string   `json:"batchcomplete"`
	Continue      Continue `json:"continue"`
	Query         Query    `json:"query"`
}

// --------------- DETALLE PARA UNA IMAGEN ---------------------------------------
// ImageInfo contiene detalles específicos de la imagen.
type ImageInfo struct {
	User                string                 `json:"user"`
	URL                 string                 `json:"url"`
	DescriptionURL      string                 `json:"descriptionurl"`
	DescriptionShortURL string                 `json:"descriptionshorturl"`
	ExtMetadata         map[string]interface{} `json:"extmetadata"`
}

// Page contiene detalles de la página relacionada con la imagen.
type PageDetail struct {
	NS                  int         `json:"ns"`
	Title               string      `json:"title"`
	Missing             string      `json:"missing"`
	Known               string      `json:"known"`
	ImageRepository     string      `json:"imagerepository"`
	ImageInfo           []ImageInfo `json:"imageinfo"`
	LicenseShortName    string      `json:"LicenseShortName"`
	UsageTerms          string      `json:"UsageTerms"`
	AttributionRequired bool        `json:"AttributionRequired"`
	LicenseUrl          string      `json:"LicenseUrl"`
	Copyrighted         bool        `json:"Copyrighted"`
	Restrictions        string      `json:"Restrictions"`
	License             string      `json:"License"`
}

// Query contiene los resultados de la consulta.
type QueryDetail struct {
	Normalized []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"normalized"`
	Pages map[string]PageDetail `json:"pages"`
}

// Continue indica los datos para la paginación o continuación.
type ContinueDetail struct {
	IiStart  string `json:"iistart"`
	Continue string `json:"continue"`
}

type WikimediaDetailAPIResponse struct {
	Continue ContinueDetail `json:"continue"`
	Query    QueryDetail    `json:"query"`
}

// Estructura para mapear el JSON parcialmente
type ResponseImageData struct {
	Query QueryImageData `json:"query"`
}

// QueryData contiene la información principal de la consulta
type QueryImageData struct {
	Pages map[string]PageImageData `json:"pages"`
}

// PageData representa una página de resultados
type PageImageData struct {
	ImageInfo []ImageInfoData `json:"imageinfo"`
}

// ImageInfoData representa la información de una imagen
type ImageInfoData struct {
	URL            string          `json:"url"`            // URL directa de la imagen
	DescriptionURL string          `json:"descriptionurl"` // URL a la descripción de la imagen
	ExtMetadata    ExtMetadataData `json:"extmetadata"`    // Metadatos adicionales
}

// ExtMetadataData contiene los metadatos de la imagen
type ExtMetadataData struct {
	ObjectName       MetadataValue `json:"ObjectName"`       // Nombre del objeto
	Artist           MetadataValue `json:"Artist"`           // Artista
	LicenseShortName MetadataValue `json:"LicenseShortName"` // Nombre abreviado de la licencia
	LicenseUrl       MetadataValue `json:"LicenseUrl"`       // URL de la licencia
}

// MetadataValue encapsula valores comunes en los metadatos
type MetadataValue struct {
	Value string `json:"value"` // Valor del metadato
}

type GetBirdImageDetailResult struct {
	URL              string        `json:"url"`              // URL directa de la imagen
	DescriptionURL   string        `json:"descriptionurl"`   // URL a la descripción de la imagen
	ObjectName       MetadataValue `json:"ObjectName"`       // Nombre del objeto
	Artist           string        `json:"Artist"`           // Artista
	LicenseShortName string        `json:"LicenseShortName"` // Nombre abreviado de la licencia
	LicenseUrl       string        `json:"LicenseUrl"`       // URL de la licencia
}

// callExternalAPI simula una llamada a una API externa.
func GetBirdImageDetail(title string, resultChan chan<- GetBirdImageDetailResult, wg *sync.WaitGroup) {
	defer wg.Done() // Asegura que el contador se disminuye una vez termine esta goroutine.

	urlBirdDetail := "https://en.wikipedia.org/w/api.php?action=query&prop=imageinfo&iiprop=extmetadata|user|url&titles=File:" + url.QueryEscape(title) + "&format=json"

	var apiResponse ResponseImageData

	err := fetchWithRetry(urlBirdDetail, &apiResponse, true)
	if err != nil {
		return
	}

	// Lista de nombres de licencias gratuitas permitidas
	freeLicenses := []string{
		"CC BY 2.0",
		"CC BY-SA 2.0",
		"CC BY 2.5",
		"CC BY-SA 2.5",
		"CC BY 3.0",
		"CC BY-SA 3.0",
		"CC BY 4.0",
		"CC BY-SA 4.0",
		"CC0 1.0",
		"CC0",
	}

	// Procesar los datos del JSON
	for _, page := range apiResponse.Query.Pages {
		for _, imageInfo := range page.ImageInfo {
			license := imageInfo.ExtMetadata.LicenseShortName.Value

			// Validar si la licencia es gratuita
			isFreeLicense := false
			for _, freeLicense := range freeLicenses {
				if license == freeLicense {
					isFreeLicense = true
					break
				}
			}

			// Si la licencia no es gratuita, ignorar este registro
			if !isFreeLicense {
				continue
			}

			// Extraer y crear el objeto
			artist := extractArtist(imageInfo.ExtMetadata.Artist.Value)
			image := GetBirdImageDetailResult{
				ObjectName:       imageInfo.ExtMetadata.ObjectName,
				URL:              imageInfo.URL,
				DescriptionURL:   imageInfo.DescriptionURL,
				Artist:           artist,
				LicenseShortName: license,
				LicenseUrl:       imageInfo.ExtMetadata.LicenseUrl.Value,
			}

			resultChan <- image // Enviar resultado al canal
		}
	}
}

// Llamada a la API de Wikipedia
func GetAllBirdImageDetails(nameBird string) ([]GetBirdImageDetailResult, error) {
	// Comprobar si el nombre científico está en el mapa y usar el nombre correcto
	if correctName, exists := latinToCommonName[nameBird]; exists {
		nameBird = correctName
	}

	// Crear un arreglo para almacenar todos los resultados
	var allImages []GetBirdImageDetailResult

	url := fmt.Sprintf("https://commons.wikimedia.org/w/api.php?action=query&format=json&list=search&srsearch=%s&srnamespace=6&utf8=&srlimit=15",
		ReplaceSpacesWithUnderscore(nameBird))

	var apiResponse WikimediaAPIResponse

	err := fetchWithRetry(url, &apiResponse, true)
	if err != nil {
		return []GetBirdImageDetailResult{}, err
	}

	// fmt.Printf("Resultados totales: %d\n", apiResponse.Query.SearchInfo.TotalHits)

	var wg sync.WaitGroup
	resultChan := make(chan GetBirdImageDetailResult, 15) // Canal para resultados

	// Iterar sobre cada título
	for _, result := range apiResponse.Query.Search {
		// Eliminamos "File:" del título
		sanitizedTitle := removePrefix(result.Title, "File:")
		//fmt.Printf("Título: %s\n", sanitizedTitle)

		wg.Add(1) // Incrementar el contador de goroutines

		go GetBirdImageDetail(sanitizedTitle, resultChan, &wg)
	}

	// Esperar a que todas las goroutines terminen
	go func() {
		wg.Wait()
		close(resultChan) // Cerrar el canal cuando todo haya terminado
	}()

	// Recoger los resultados del canal
	for image := range resultChan {
		allImages = append(allImages, image)
	}

	return allImages, nil
}

// removePrefix elimina un prefijo dado de un string.
func removePrefix(text, prefix string) string {
	if len(text) >= len(prefix) && text[:len(prefix)] == prefix {
		return text[len(prefix):]
	}
	return text
}

func extractArtist(value string) string {
	// Extraer el contenido dentro de las etiquetas <a>
	start := strings.Index(value, ">") + 1
	end := strings.LastIndex(value, "<")
	if start > 0 && end > start {
		return value[start:end]
	}
	return ""
}
