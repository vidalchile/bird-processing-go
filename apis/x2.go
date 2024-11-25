package apis

import (
	"fmt"
	"math/rand"
	"project/models"
	"time"
)

// SimulateAPIX2 simula la API X2 devolviendo una descripción
func SimulateAPIX2(recordID int) (*models.Description, error) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Simula delay
	return &models.Description{Description: "Descripción para el registro " + fmt.Sprint(recordID)}, nil
}
