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

func buildURL(base string, pathSegments ...string) string {
	// Usa strings.Join para combinar los segmentos con una sola barra '/'
	return fmt.Sprintf("%s/%s", base, strings.Join(pathSegments, "/"))
}

// reemplaza los espacios en blanco de una cadena con "_"
func ReplaceSpacesWithUnderscore(input string) string {
	return strings.ReplaceAll(input, " ", "_")
}

// Llamada a la API de Wikipedia
func GetBirdExtract(nameBird string) (string, error) {
	url := fmt.Sprintf(buildURL("https://es.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts|pageimages&explaintext=true&ppprop=original&origin=*"),
		ReplaceSpacesWithUnderscore(nameBird))

	var apiResponse WikipediaAPIResponse

	err := fetchAndParseJSON(url, &apiResponse)
	if err != nil {
		return "", err
	}

	// Obtener la descripción de la primera página encontrada
	for _, page := range apiResponse.Query.Pages {
		if page.Extract != "" {
			return page.Extract, nil
		}
	}

	return "", errors.New("no se encontró una descripción para el pájaro en Wikipedia")
}
