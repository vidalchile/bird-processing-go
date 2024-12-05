package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Función auxiliar para hacer la solicitud HTTP y deserializar la respuesta JSON.
func fetchAndParseJSON(url string, response interface{}) error {
	// Asegurarse de que no haya barras extra en la URL
	// url = strings.ReplaceAll(url, "//", "/")

	log.Println(">>> url: ", url)

	// Crear un cliente HTTP para agregar las cabeceras
	client := &http.Client{}

	// Crear una solicitud GET con cabeceras personalizadas
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error al crear la solicitud HTTP: %v", err)
	}

	// Agregar cabeceras para parecer un navegador real
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,es;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Realizar la solicitud
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al hacer la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado HTTP
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta HTTP: código %d, URL: %s", resp.StatusCode, url)
	}

	// Leer la respuesta usando io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Deserializar la respuesta JSON
	err = json.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("error al deserializar la respuesta JSON: %v", err)
	}

	// Agregar un pequeño retraso de 500ms para evitar hacer demasiadas solicitudes de golpe
	time.Sleep(500 * time.Millisecond)

	return nil
}

// fetchWithRetry intenta realizar la solicitud varias veces si ocurre un error temporal.
func fetchWithRetry(url string, response interface{}) error {
	const maxRetries = 3 // Número máximo de reintentos
	var err error

	for i := 0; i < maxRetries; i++ {
		// Intentar obtener los datos usando la función fetchAndParseJSON
		err = fetchAndParseJSON(url, response)
		if err == nil {
			// Si la solicitud fue exitosa, salimos de la función
			return nil
		}

		// Si no fue exitosa, registramos el error y esperamos antes de reintentar
		log.Printf("Intento %d fallido: %v. Reintentando...", i+1, err)

		// Espera de 2 segundos antes de intentar de nuevo
		time.Sleep(2 * time.Second)
	}

	// Si fallan todos los intentos, devolvemos el error final
	return fmt.Errorf("no se pudo completar la solicitud después de %d intentos: %v", maxRetries, err)
}
