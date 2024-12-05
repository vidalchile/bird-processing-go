package models

// Bird representa un registro de un ave
type Bird struct {
	UID    string     `json:"uid"`
	Name   BirdName   `json:"name"`
	Images BirdImages `json:"images"`
	Links  BirdLinks  `json:"_links"`
	Sort   int        `json:"sort"`
}

// BirdName representa los nombres de un ave
type BirdName struct {
	Spanish string `json:"spanish"`
	English string `json:"english"`
	Latin   string `json:"latin"`
}

// BirdImages contiene URLs de imágenes del ave
type BirdImages struct {
	Main  string `json:"main"`
	Full  string `json:"full"`
	Thumb string `json:"thumb"`
}

// BirdLinks contiene enlaces relacionados con el ave
type BirdLinks struct {
	Self   string `json:"self"`
	Parent string `json:"parent"`
}

// Image representa una imagen devuelta por la API X3
type Image struct {
	URL string `json:"url"`
}
