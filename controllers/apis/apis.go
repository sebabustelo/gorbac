package apis

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"

	responses "api-rbac/helpers"
	"api-rbac/models"
)

var (
	publicKey *rsa.PublicKey
)

// Index retorna un listado de usuarios
func Index(w http.ResponseWriter, r *http.Request) {

	api := models.Api{}
	apis, err := api.FindAll()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Error al generar el json"))
		return
	}

	responses.JSON(w, http.StatusOK, apis)

}

// Add ...
func Add(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userAdd, err := user.Create(user)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, userAdd)
}

// Edit permite realizar la modificación de un usuario, recibe el id del mismo por post
func Edit(w http.ResponseWriter, r *http.Request) {

	// var application models.Application
	// err := json.NewDecoder(r.Body).Decode(&application)

	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	// applicationUpdate, err := application.Update(application)

	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	// responses.JSON(w, http.StatusOK, applicationUpdate)
}

// GetByID ...
func GetByID(w http.ResponseWriter, r *http.Request) {

	// idApplication := chi.URLParam(r, "id")

	// id, err := strconv.Atoi(idApplication)
	// var application models.Application
	// result, err := application.GetByID(id)

	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	// responses.JSON(w, http.StatusOK, result)
}
