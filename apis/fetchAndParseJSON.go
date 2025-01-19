package apis

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// CustomClient es un cliente HTTP con configuraciones extendidas para las solicitudes
var customClient = &http.Client{
	Transport: &http.Transport{
		// Configuración de TLS para garantizar que no se omita la verificación de certificados
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // No omitir verificación de certificados, garantizar seguridad
		},
		// Configuración de red para manejar las conexiones de red
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // Tiempo de espera para la conexión antes de dar error
			KeepAlive: 30 * time.Second, // Mantener la conexión abierta por 10 segundos
		}).DialContext,
	},
}

// fetchAndParseJSON realiza una solicitud HTTP y deserializa la respuesta JSON en el objeto `response` proporcionado.
func fetchAndParseJSON(url string, response interface{}, useCustomClient bool) error {
	// Crear un cliente HTTP estándar (por defecto)
	client := &http.Client{}

	// Si la variable `useCustomClient` es true, utilizamos el `customClient` con configuraciones extendidas
	if useCustomClient {
		client = customClient
	}

	// Crear una solicitud HTTP de tipo GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// Si ocurre un error al crear la solicitud, se devuelve un error detallado
		return fmt.Errorf("error al crear la solicitud HTTP: %v", err)
	}

	// Agregar cabeceras personalizadas a la solicitud para simular un navegador real
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*") // Espera recibir JSON
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,es;q=0.8")  // Idiomas aceptables
	req.Header.Set("Connection", "keep-alive")                    // Mantener la conexión viva
	req.Header.Set("Upgrade-Insecure-Requests", "1")              // Solicitar contenido seguro

	// Realizar la solicitud HTTP usando el cliente configurado
	resp, err := client.Do(req)
	if err != nil {
		// Si ocurre un error en la solicitud (por ejemplo, timeout o problema de red), se devuelve un error
		return fmt.Errorf("error al hacer la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close() // Asegurar que el cuerpo de la respuesta se cierre al final

	// Verificar que el código de estado HTTP sea 200 OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error HTTP %d: %s", resp.StatusCode, body)
		// Si el código de estado es diferente de 200, devolver un error con el código y URL
		return fmt.Errorf("error en la respuesta HTTP: código %d, URL: %s", resp.StatusCode, url)
	}

	// Leer todo el cuerpo de la respuesta HTTP
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Si ocurre un error al leer la respuesta, se devuelve un error detallado
		return fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Deserializar el cuerpo de la respuesta JSON en el objeto `response`
	err = json.Unmarshal(body, response)
	if err != nil {
		// Si ocurre un error al deserializar el JSON, se devuelve un error
		return fmt.Errorf("error al deserializar la respuesta JSON: %v", err)
	}

	// Si todo ha ido bien, la función retorna nil (sin errores)
	return nil
}

// fetchWithRetry intenta realizar la solicitud varias veces si ocurre un error temporal, hasta un máximo de `maxRetries` intentos.
func fetchWithRetry(url string, response interface{}, useCustomClient bool) error {
	const maxRetries = 3 // Número máximo de intentos en caso de error
	var err error

	// Intentar hacer la solicitud hasta `maxRetries` veces
	for i := 0; i < maxRetries; i++ {
		// Intentamos obtener los datos usando la función fetchAndParseJSON
		err = fetchAndParseJSON(url, response, useCustomClient)
		if err == nil {
			// Si la solicitud fue exitosa (sin errores), retornamos nil
			return nil
		}

		// Si la solicitud falla, registramos el error y esperamos antes de reintentar
		log.Printf("Intento %d fallido: %v. Reintentando...", i+1, err)

		// Espera de 2 segundos antes de intentar nuevamente
		// time.Sleep(2 * time.Second)
	}

	// Si todos los intentos fallaron, devolvemos un error final indicando que no se pudo completar la solicitud
	return fmt.Errorf("no se pudo completar la solicitud después de %d intentos: %v", maxRetries, err)
}
