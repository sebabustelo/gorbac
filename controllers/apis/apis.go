package apis

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"api-rbac/db"
	responses "api-rbac/helpers"
	"api-rbac/models"

	"github.com/go-chi/chi"
)

var (
	publicKey *rsa.PublicKey
)

// SetRoutes establece las rutas dinámicamente desde main.go
func SetRoutes(routes []string) {
	allRoutes = routes
}

// Lista de rutas del router (se establece dinámicamente desde main.go)
var allRoutes = []string{
	"/users/add",
	"/users/edit",
	"/users/delete/{id}",
	"/users/{id}",
	"/users/index",
	"/roles",
	"/roles/add",
	"/roles/{id}",
	"/products/add",
	"/products/{id}",
	"/apis",
	"/apis/add",
	"/apis/{id}",
	"/roles/{id}/apis",
	"/login",
	"/roles/permissions/{id}/apis",
	"/refresh",
	"/products",
	"/google-login",
	"/health",
	"/auth/check",
	// Orders routes
	"/orders",
	"/orders/{id}",
	"/orders/{id}/status",
	"/orders/user/{user_id}",
}

// Index retorna un listado de todas las APIs
func Index(w http.ResponseWriter, r *http.Request) {
	api := models.Api{}
	apis, err := api.FindAll()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("error al obtener las APIs"))
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

	if api.Tipo == "" {
		api.Tipo = "GET"
	}

	// Validar campos requeridos
	if api.Endpoint == "" {
		responses.ERROR(w, http.StatusBadRequest, errors.New("el endpoint es requerido"))
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

	if apiUpdate.Tipo == "" {
		apiUpdate.Tipo = "GET"
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
	existingApi.Tipo = apiUpdate.Tipo

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

// GET /apis/sync: compara rutas del router con la tabla apis
func SyncApis(w http.ResponseWriter, r *http.Request) {
	db := db.Instance()
	var apis []models.Api
	db.Find(&apis)

	dbEndpoints := make(map[string]bool)
	for _, api := range apis {
		dbEndpoints[api.Endpoint] = true
	}

	missing := []string{}
	existing := []string{}
	for _, route := range allRoutes {
		if dbEndpoints[route] {
			existing = append(existing, route)
		} else {
			missing = append(missing, route)
		}
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"missing":  missing,
		"existing": existing,
	})
}

// POST /apis/sync: agrega endpoints faltantes a la tabla apis
func AddMissingApis(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Endpoints []string `json:"endpoints"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	db := db.Instance()
	added := []string{}
	for _, ep := range payload.Endpoints {
		var count int64
		db.Model(&models.Api{}).Where("endpoint = ?", ep).Count(&count)
		if count == 0 {
			api := models.Api{Endpoint: ep, Description: ""}
			db.Create(&api)
			added = append(added, ep)
		}
	}
	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"added": added,
	})
}
