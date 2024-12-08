package apis

import (
	"errors"
	"fmt"
	"strings"
)

// Structs para manejar la respuesta de la API de Wikipedia
type WikipediaAPIResponse struct {
	Query QueryResult `json:"query"`
}

type QueryResult struct {
	Pages map[string]Page `json:"pages"`
}

type Page struct {
	Title     string `json:"title"`
	Extract   string `json:"extract"`
	Thumbnail struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
}

// reemplaza los espacios en blanco de una cadena con "_"
func ReplaceSpacesWithUnderscore(input string) string {
	return strings.ReplaceAll(input, " ", "_")
}

// Llamada a la API de Wikipedia
func GetBirdExtract(nameBird string) (string, error) {
	url := fmt.Sprintf("https://es.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts|pageimages&explaintext=true&ppprop=original&origin=*",
		ReplaceSpacesWithUnderscore(nameBird))

	var apiResponse WikipediaAPIResponse

	err := fetchWithRetry(url, &apiResponse, true)
	if err != nil {
		return "", err
	}

	// Obtener la descripción de la primera página encontrada
	for _, page := range apiResponse.Query.Pages {
		if page.Extract != "" {
			return page.Extract, nil
		}
	}

	return "", errors.New("No se encontró una descripción para el pájaro " + nameBird + " en Wikipedia")
}
