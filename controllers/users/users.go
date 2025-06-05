package users

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	responses "api-rbac/helpers"
	"api-rbac/models"
)

var (
	publicKey  *rsa.PublicKey
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// validateUser realiza la validación de los campos del usuario
func validateUser(user models.User, isUpdate bool) error {
	if strings.TrimSpace(user.Name) == "" {
		return fmt.Errorf("el nombre es requerido")
	}

	if len(user.Name) < 2 || len(user.Name) > 50 {
		return fmt.Errorf("el nombre debe tener entre 2 y 50 caracteres")
	}

	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}

	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("el formato del email no es válido")
	}

	if !isUpdate && strings.TrimSpace(user.Password) == "" {
		return fmt.Errorf("la contraseña es requerida")
	}

	if user.Password != "" && (len(user.Password) < 6 || len(user.Password) > 50) {
		return fmt.Errorf("la contraseña debe tener entre 6 y 50 caracteres")
	}

	return nil
}

// Index retorna un listado de usuarios con paginación y filtros
func Index(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de paginación
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Obtener parámetros de búsqueda
	search := r.URL.Query().Get("search")

	usr := models.User{}
	users, total, err := usr.FindAllWithPagination(page, limit, search)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Calcular metadata de paginación
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := map[string]interface{}{
		"data": users,
		"pagination": map[string]interface{}{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	}

	responses.JSON(w, http.StatusOK, response)
}

// Add crea un nuevo usuario en el sistema
func Add(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Validación de campos
	if err := validateUser(user, false); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userAdd, err := user.Create(user)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, userAdd)
}

// Edit permite realizar la modificación de un usuario, recibe el id del mismo por post
func Edit(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Validación de ID
	if user.ID == 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("el ID del usuario es requerido"))
		return
	}

	// Validación de campos
	if err := validateUser(user, true); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userUpdate, err := user.Update(user)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, userUpdate)
}

// Delete elimina un usuario del sistema por su ID
func Delete(w http.ResponseWriter, r *http.Request) {

	idUser := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idUser)

	if err != nil {
		responses.ERROR(w, http.StatusNotAcceptable, err)
		return
	}

	var userDelete models.User
	result, err := userDelete.Delete(id)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

// GetByID obtiene un usuario específico por su ID
func GetByID(w http.ResponseWriter, r *http.Request) {

	idUser := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idUser)

	if err != nil {
		responses.ERROR(w, http.StatusNotAcceptable, err)
		return
	}

	var user models.User
	result, err := user.GetByID(id)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, result)
}
