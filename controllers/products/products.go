package products

import (
	"encoding/json"
	"net/http"

	"api-rbac/models"
	"api-rbac/responses"
)

func Index(w http.ResponseWriter, r *http.Request) {

	product := models.Product{}
	products, err := product.FindAll()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, products)
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
