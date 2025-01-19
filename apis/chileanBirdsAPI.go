package apis

import "project/models"

type ResultBirdDetail struct {
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

// Llamada a la API de Chilean Birds para obtener listado de aves (TEST 2 aves)
func GetBirdsTest() ([]models.Bird, error) {
	records := []models.Bird{
		{
			UID: "76-buteo-albigula",
			Name: models.BirdName{
				Spanish: "Martín Pescador",
				English: "Ringed Kingfisher",
				Latin:   "Ceryle torquata",
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
	}

	return records, nil
}

// Llamada a la API de Chilean Birds para obtener listado de aves
func GetBirds() ([]models.Bird, error) {
	url := "https://aves.ninjas.cl/api/birds"

	var records []models.Bird

	err := fetchWithRetry(url, &records, true)
	if err != nil {
		return records, err
	}

	return records, nil
}

// Llamada a la API de Chilean Birds para obtener el detalle de una ave
func GetBirdDetail(url string) (ResultBirdDetail, error) {
	var apiResponse ResultBirdDetail

	err := fetchWithRetry(url, &apiResponse, true)
	if err != nil {
		return ResultBirdDetail{}, err
	}

	return apiResponse, nil
}
