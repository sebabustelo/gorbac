package products

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"api-rbac/models"
	"api-rbac/responses"

	"github.com/go-chi/chi"
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
		fmt.Printf("Error decodificando JSON: %v\n", err)
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Validar que los campos requeridos estén presentes
	if product.Name == "" {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("nombre del producto es requerido"))
		return
	}

	if product.Price <= 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("precio debe ser mayor a 0"))
		return
	}

	if product.Stock < 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("stock no puede ser negativo"))
		return
	}

	// Si no se especifica categoría, usar categoría por defecto (ID 1)
	if product.CategoryID == 0 {
		product.CategoryID = 1
	}

	fmt.Printf("Creando producto: %+v\n", product)

	productAdd, err := product.Create()

	if err != nil {
		fmt.Printf("Error creando producto: %v\n", err)
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Recargar el producto con la categoría para obtener el nombre
	productWithCategory, err := productAdd.GetByID(int(productAdd.ID))
	if err != nil {
		fmt.Printf("Error obteniendo producto con categoría: %v\n", err)
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Convertir a la estructura de respuesta del frontend
	productResponse := ProductResponse{
		ID:          strconv.FormatUint(uint64(productWithCategory.ID), 10),
		Nombre:      productWithCategory.Name,
		Descripcion: productWithCategory.Description,
		Precio:      productWithCategory.Price,
		Stock:       productWithCategory.Stock,
		Categoria:   productWithCategory.Category.Name,
		Imagen:      productWithCategory.Image,
	}

	fmt.Printf("Producto creado exitosamente: %+v\n", productResponse)
	responses.JSON(w, http.StatusOK, productResponse)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del producto de la URL
	productID := chi.URLParam(r, "id")
	if productID == "" {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de producto requerido"))
		return
	}

	// Convertir el ID a uint
	id, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de producto inválido"))
		return
	}

	// Decodificar el cuerpo de la petición
	var productUpdate models.Product
	if err := json.NewDecoder(r.Body).Decode(&productUpdate); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Obtener el producto existente
	product := models.Product{}
	existingProduct, err := product.GetByID(int(id))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("producto no encontrado"))
		return
	}

	// Actualizar los campos del producto
	existingProduct.Name = productUpdate.Name
	existingProduct.Description = productUpdate.Description
	existingProduct.Price = productUpdate.Price
	existingProduct.Stock = productUpdate.Stock
	existingProduct.CategoryID = productUpdate.CategoryID
	existingProduct.Image = productUpdate.Image

	// Guardar los cambios
	updatedProduct, err := existingProduct.Update()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Convertir a la estructura de respuesta del frontend
	productResponse := ProductResponse{
		ID:          strconv.FormatUint(uint64(updatedProduct.ID), 10),
		Nombre:      updatedProduct.Name,
		Descripcion: updatedProduct.Description,
		Precio:      updatedProduct.Price,
		Stock:       updatedProduct.Stock,
		Categoria:   updatedProduct.Category.Name,
		Imagen:      updatedProduct.Image,
	}

	responses.JSON(w, http.StatusOK, productResponse)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del producto de la URL
	productID := chi.URLParam(r, "id")
	if productID == "" {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de producto requerido"))
		return
	}

	// Convertir el ID a uint
	id, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de producto inválido"))
		return
	}

	// Obtener el producto existente
	product := models.Product{}
	existingProduct, err := product.GetByID(int(id))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("producto no encontrado"))
		return
	}

	// Eliminar el producto
	err = existingProduct.Delete()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "Producto eliminado correctamente"})
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de producto requerido"))
		return
	}

	id, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("ID de producto inválido"))
		return
	}

	product := models.Product{}
	existingProduct, err := product.GetByID(int(id))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("producto no encontrado"))
		return
	}

	productResponse := ProductResponse{
		ID:          strconv.FormatUint(uint64(existingProduct.ID), 10),
		Nombre:      existingProduct.Name,
		Descripcion: existingProduct.Description,
		Precio:      existingProduct.Price,
		Stock:       existingProduct.Stock,
		Categoria:   existingProduct.Category.Name,
		Imagen:      existingProduct.Image,
	}

	responses.JSON(w, http.StatusOK, productResponse)
}
