package models

//ResponseToken representa la estructura para devolver el token codificado en base64
type ResponseToken struct {
	Token string `json:"token"`
}
