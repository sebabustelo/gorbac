package orders

import (
	"api-rbac/models"
	"api-rbac/responses"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type CreateOrderRequest struct {
	UserID uint               `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Calcular total y crear items
	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range req.Items {
		// Obtener producto para verificar stock y precio
		product := models.Product{}
		foundProduct, err := product.GetByID(int(item.ProductID))
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		// Verificar stock
		if foundProduct.Stock < item.Quantity {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		// Calcular subtotal
		subtotal := foundProduct.Price * float64(item.Quantity)
		totalAmount += subtotal

		// Crear item
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: foundProduct.Price,
			Subtotal:  subtotal,
		}
		orderItems = append(orderItems, orderItem)

		// Actualizar stock
		foundProduct.Stock -= item.Quantity
		foundProduct.Update()
	}

	// Crear orden
	order := models.Order{
		UserID:      req.UserID,
		TotalAmount: totalAmount,
		Status:      "pending",
		Items:       orderItems,
	}

	createdOrder, err := order.Create()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, createdOrder)
}

func Index(w http.ResponseWriter, r *http.Request) {
	order := models.Order{}
	orders, err := order.FindAll()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, orders)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	order := models.Order{}
	foundOrder, err := order.GetByID(id)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	responses.JSON(w, http.StatusOK, foundOrder)
}

func GetByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	order := models.Order{}
	orders, err := order.GetByUser(userID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, orders)
}

func UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	order := models.Order{}
	if err := order.UpdateStatus(id, req.Status); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "Estado actualizado correctamente"})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	order := models.Order{}
	if err := order.Delete(id); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "Orden eliminada correctamente"})
}
