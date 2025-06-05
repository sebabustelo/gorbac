package roles

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"api-rbac/models"

	responses "api-rbac/helpers"

	"github.com/go-chi/chi"
)

var (
	publicKey *rsa.PublicKey
)

// Index retorna un listado de usuarios
func Index(w http.ResponseWriter, r *http.Request) {

	role := models.Role{}
	roles, err := role.FindAll()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Error al generar el json"))
		return
	}

	responses.JSON(w, http.StatusOK, roles)

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

// GetByID ...
func GetByID(w http.ResponseWriter, r *http.Request) {

	idRole := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idRole)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("ID inválido"))
		return
	}
	var role models.Role
	result, err := role.GetByID(id)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

func GetApisByRole(w http.ResponseWriter, r *http.Request) {
	idRole := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idRole)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("ID inválido"))
		return
	}

	role := models.Role{}
	result, err := role.GetApisByRole(id)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, result)
}
