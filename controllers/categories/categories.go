package categories

import (
	"api-rbac/models"
	"api-rbac/responses"
	"encoding/json"
	"net/http"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	categories, err := category.FindAll()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	// Opcional: transformar la respuesta si no quieres enviar los productos
	for i := range categories {
		categories[i].Products = nil
	}
	json.NewEncoder(w).Encode(categories)
}
