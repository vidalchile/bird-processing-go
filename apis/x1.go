package apis

import (
	"fmt"

	"project/models"
)

// SimulateAPIX1 simula la API X1 devolviendo 300 registros
func SimulateAPIX1() ([]models.Record, error) {
	records := make([]models.Record, 300)
	for i := 0; i < 300; i++ {
		records[i] = models.Record{ID: i + 1, Name: fmt.Sprintf("Record-%d", i+1)}
	}
	return records, nil
}
