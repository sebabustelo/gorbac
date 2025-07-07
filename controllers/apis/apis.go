package apis

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	responses "api-rbac/helpers"
	"api-rbac/models"

	"github.com/go-chi/chi"
)

var (
	publicKey *rsa.PublicKey
)

// Index retorna un listado de todas las APIs
func Index(w http.ResponseWriter, r *http.Request) {
	api := models.Api{}
	apis, err := api.FindAll()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Error al obtener las APIs"))
		return
	}

	responses.JSON(w, http.StatusOK, apis)
}

// Add crea una nueva API
func Add(w http.ResponseWriter, r *http.Request) {
	var api models.Api
	err := json.NewDecoder(r.Body).Decode(&api)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Validar campos requeridos
	if api.Endpoint == "" {
		responses.ERROR(w, http.StatusBadRequest, errors.New("El endpoint es requerido"))
		return
	}

	apiAdd, err := api.Create()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, apiAdd)
}

// Edit permite modificar una API
func Edit(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID de la API de la URL
	apiID := chi.URLParam(r, "id")
	if apiID == "" {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de API requerido"))
		return
	}

	// Convertir el ID a uint
	id, err := strconv.ParseUint(apiID, 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de API inválido"))
		return
	}

	// Decodificar el cuerpo de la petición
	var apiUpdate models.Api
	if err := json.NewDecoder(r.Body).Decode(&apiUpdate); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Obtener la API existente
	api := models.Api{}
	existingApi, err := api.GetByID(int(id))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("API no encontrada"))
		return
	}

	// Actualizar los campos de la API
	existingApi.Endpoint = apiUpdate.Endpoint
	existingApi.Description = apiUpdate.Description
	existingApi.Hidden = apiUpdate.Hidden
	existingApi.Public = apiUpdate.Public

	// Guardar los cambios
	updatedApi, err := existingApi.Update()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, updatedApi)
}

// Delete elimina una API
func Delete(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID de la API de la URL
	apiID := chi.URLParam(r, "id")
	if apiID == "" {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de API requerido"))
		return
	}

	// Convertir el ID a uint
	id, err := strconv.ParseUint(apiID, 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de API inválido"))
		return
	}

	// Obtener la API existente
	api := models.Api{}
	existingApi, err := api.GetByID(int(id))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("API no encontrada"))
		return
	}

	// Eliminar la API
	err = existingApi.Delete()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "API eliminada correctamente"})
}

// GetByID obtiene una API por su ID
func GetByID(w http.ResponseWriter, r *http.Request) {
	idAPI := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idAPI)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("ID inválido"))
		return
	}

	var api models.Api
	result, err := api.GetByID(id)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, result)
}
