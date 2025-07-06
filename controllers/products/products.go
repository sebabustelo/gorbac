package products

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-rbac/models"
	"api-rbac/responses"
)

// ProductResponse representa la estructura que espera el frontend
type ProductResponse struct {
	ID          string  `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	Categoria   string  `json:"categoria"`
	Imagen      string  `json:"imagen"`
}

func Index(w http.ResponseWriter, r *http.Request) {

	product := models.Product{}
	products, err := product.FindAll()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Convertir los productos al formato que espera el frontend
	var productResponses []ProductResponse
	for _, p := range products {
		productResponse := ProductResponse{
			ID:          strconv.FormatUint(uint64(p.ID), 10),
			Nombre:      p.Name,
			Descripcion: p.Description,
			Precio:      p.Price,
			Stock:       p.Stock,
			Categoria:   p.Category.Name,
			Imagen:      p.Image,
		}
		productResponses = append(productResponses, productResponse)
	}

	responses.JSON(w, http.StatusOK, productResponses)
}

func Add(w http.ResponseWriter, r *http.Request) {

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	productAdd, err := product.Create()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, productAdd)
}

func Edit(w http.ResponseWriter, r *http.Request) {

	var product models.Product
	json.NewDecoder(r.Body).Decode(&product)

}

func Delete(w http.ResponseWriter, r *http.Request) {

	var product models.Product
	json.NewDecoder(r.Body).Decode(&product)

}
