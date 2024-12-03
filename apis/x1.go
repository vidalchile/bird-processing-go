package apis

import (
	"project/models"
)

// SimulateAPIX1 simula la API X1 devolviendo datos de aves
func SimulateAPIX1() ([]models.Bird, error) {
	records := []models.Bird{
		{
			UID: "76-buteo-albigula",
			Name: models.BirdName{
				Spanish: "Aguilucho Chico",
				English: "White-throated Hawk",
				Latin:   "Buteo albigula",
			},
			Images: models.BirdImages{
				Main:  "https://aves.ninjas.cl/api/site/assets/files/3099/17082018024245aguilucho_chico_tomas_rivas_web.200x0.jpg",
				Full:  "https://aves.ninjas.cl/api/site/assets/files/3099/17082018024245aguilucho_chico_tomas_rivas_web.jpg",
				Thumb: "https://aves.ninjas.cl/api/site/assets/files/3099/17082018024245aguilucho_chico_tomas_rivas_web.200x0.jpg",
			},
			Links: models.BirdLinks{
				Self:   "https://aves.ninjas.cl/api/birds/76-buteo-albigula",
				Parent: "https://aves.ninjas.cl/api/birds",
			},
			Sort: 0,
		},
		{
			UID: "46-lophonetta-specularioides",
			Name: models.BirdName{
				Spanish: "Pato Juarjual",
				English: "Crested Duck",
				Latin:   "Lophonetta specularioides",
			},
			Images: models.BirdImages{
				Main:  "https://aves.ninjas.cl/api/site/assets/files/3102/18082018072023pato_juarjual_pedro_valencia_web.200x0.jpg",
				Full:  "https://aves.ninjas.cl/api/site/assets/files/3102/18082018072023pato_juarjual_pedro_valencia_web.jpg",
				Thumb: "https://aves.ninjas.cl/api/site/assets/files/3102/18082018072023pato_juarjual_pedro_valencia_web.200x0.jpg",
			},
			Links: models.BirdLinks{
				Self:   "https://aves.ninjas.cl/api/birds/46-lophonetta-specularioides",
				Parent: "https://aves.ninjas.cl/api/birds",
			},
			Sort: 1,
		},
	}

	return records, nil
}
