package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type QueryResultDetail struct {
	Uid        string `json:"uid"`
	Map        any    `json:"map"`
	Iucn       any    `json:"iucn"`
	Migration  bool   `json:"migration"`
	Dimorphism bool   `json:"dimorphism"`
	Size       string `json:"size"`
	Order      string `json:"order"`
	Species    string `json:"species"`
	Images     any    `json:"images"`
	Audio      any    `json:"audio"`
}

// Llamada a la API de Wikipedia
func CallDetailBirdAPI(url string) (QueryResultDetail, error) {
	// Hacer la solicitud HTTP
	resp, err := http.Get(url)
	if err != nil {
		return QueryResultDetail{}, fmt.Errorf("error al hacer la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta usando io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return QueryResultDetail{}, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Deserializar la respuesta JSON
	var apiResponse QueryResultDetail
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return QueryResultDetail{}, fmt.Errorf("error al deserializar la respuesta JSON: %v", err)
	}

	return apiResponse, nil
}
