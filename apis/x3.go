package apis

import (
	"fmt"
	"math/rand"
	"project/models"
	"time"
)

// SimulateAPIX3 simula la API X3 devolviendo imágenes
func SimulateAPIX3(recordID int) ([]models.Image, error) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Simula delay
	return []models.Image{
		{URL: "https://example.com/record" + fmt.Sprint(recordID) + "/image1.jpg"},
		{URL: "https://example.com/record" + fmt.Sprint(recordID) + "/image2.jpg"},
	}, nil
}
