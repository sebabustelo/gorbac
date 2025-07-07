package roles

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"api-rbac/db"
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

// UpdateApis actualiza las APIs permitidas de un rol
func UpdateApis(w http.ResponseWriter, r *http.Request) {
	idRole := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idRole)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("ID inválido"))
		return
	}

	var payload struct {
		Apis []int `json:"apis"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db := db.Instance()
	var role models.Role
	if err := db.Preload("Apis").First(&role, id).Error; err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Rol no encontrado"))
		return
	}

	// Obtener las APIs a asociar
	var apis []models.Api
	if len(payload.Apis) > 0 {
		if err := db.Where("id IN (?)", payload.Apis).Find(&apis).Error; err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
	}

	// Actualizar la relación
	if err := db.Model(&role).Association("Apis").Replace(apis).Error; err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "APIs actualizadas correctamente",
		"role":    role.ID,
		"apis":    payload.Apis,
	})
}
