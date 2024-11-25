package models

// Record representa los registros devueltos por la API X1
type Record struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Description representa la descripción devuelta por la API X2
type Description struct {
	Description string `json:"description"`
}

// Image representa una imagen devuelta por la API X3
type Image struct {
	URL string `json:"url"`
}
