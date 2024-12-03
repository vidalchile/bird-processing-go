package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Structs para manejar la respuesta de la API de Wikipedia
type Page struct {
	Title     string `json:"title"`
	Extract   string `json:"extract"`
	Thumbnail struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
}

type QueryResult struct {
	Pages map[string]Page `json:"pages"`
}

type WikipediaAPIResponse struct {
	Query QueryResult `json:"query"`
}

// ReplaceSpacesWithUnderscore reemplaza los espacios en blanco de una cadena con "_"
func ReplaceSpacesWithUnderscore(input string) string {
	return strings.ReplaceAll(input, " ", "_")
}

// Llamada a la API de Wikipedia
func CallWikipediaAPI(nameBird string) (string, error) {
	// URL de la API de Wikipedia
	url := "https://es.wikipedia.org/w/api.php?action=query&format=json&titles=" + ReplaceSpacesWithUnderscore(nameBird) + "&prop=extracts|pageimages&explaintext=true&ppprop=original&origin=*"

	// Hacer la solicitud HTTP
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error al hacer la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta usando io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Deserializar la respuesta JSON
	var apiResponse WikipediaAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("error al deserializar la respuesta JSON: %v", err)
	}

	// Obtener la página de interés
	page := apiResponse.Query.Pages
	var description string

	for _, p := range page {
		description = p.Extract
	}

	return description, nil
}
