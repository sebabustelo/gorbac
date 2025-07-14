package orders

import (
	"api-rbac/authentication"
	"api-rbac/db"
	"api-rbac/models"
	"api-rbac/responses"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

type CreateOrderRequest struct {
	Items    []OrderItemRequest `json:"items"`
	Payment  PaymentRequest     `json:"payment"`
	Shipping ShippingRequest    `json:"shipping"`
	Notes    string             `json:"notes"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type PaymentRequest struct {
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
}

type ShippingRequest struct {
	RecipientName      string `json:"recipient_name"`
	RecipientEmail     string `json:"recipient_email"`
	RecipientPhone     string `json:"recipient_phone"`
	ShippingAddress    string `json:"shipping_address"`
	ShippingCity       string `json:"shipping_city"`
	ShippingPostalCode string `json:"shipping_postal_code"`
	ShippingProvince   string `json:"shipping_province"`
	ShippingMethod     string `json:"shipping_method"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status"`
	TransactionID string `json:"transaction_id,omitempty"`
}

// generateOrderNumber genera un número de pedido único
func generateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%04d", now.Format("20060102"), now.Nanosecond()%10000)
}

func Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Create: Starting order creation")

	// Obtener el claim del contexto
	claim, ok := r.Context().Value(authentication.UserContextKey).(models.Claim)
	if !ok {
		log.Printf("Create: No claim found in context")
		responses.ERROR(w, http.StatusUnauthorized, nil)
		return
	}

	log.Printf("Create: User ID %d authenticated successfully", claim.User.ID)

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Create: JSON decode error: %v", err)
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Create: Request decoded successfully. Items count: %d", len(req.Items))
	log.Printf("Create: Payment method: %s, Amount: %.2f", req.Payment.PaymentMethod, req.Payment.Amount)
	log.Printf("Create: Shipping recipient: %s, Email: %s", req.Shipping.RecipientName, req.Shipping.RecipientEmail)

	// Validar que hay items en el pedido
	if len(req.Items) == 0 {
		log.Printf("Create: No items in order")
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("el pedido debe contener al menos un item"))
		return
	}

	log.Printf("Create: Starting transaction")
	// Iniciar transacción
	tx := db.Instance().Begin()
	if tx.Error != nil {
		log.Printf("Create: Transaction start error: %v", tx.Error)
		responses.ERROR(w, http.StatusInternalServerError, tx.Error)
		return
	}

	// Calcular totales
	var subtotal float64
	var orderItems []models.OrderItem

	// Procesar cada item del pedido
	for i, itemReq := range req.Items {
		log.Printf("Create: Processing item %d: ProductID=%d, Quantity=%d", i, itemReq.ProductID, itemReq.Quantity)

		// Obtener el producto
		var product models.Product
		if err := tx.Where("id = ?", itemReq.ProductID).First(&product).Error; err != nil {
			log.Printf("Create: Product not found: %d, error: %v", itemReq.ProductID, err)
			tx.Rollback()
			responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("producto no encontrado: %d", itemReq.ProductID))
			return
		}

		log.Printf("Create: Product found: %s, Price: %.2f", product.Name, product.Price)

		// Calcular precios
		totalPrice := float64(itemReq.Quantity) * product.Price
		subtotal += totalPrice

		// Crear order item
		orderItem := models.OrderItem{
			ProductID:    int(itemReq.ProductID),
			ProductName:  product.Name,
			ProductImage: product.Image,
			Quantity:     itemReq.Quantity,
			UnitPrice:    product.Price,
			TotalPrice:   totalPrice,
		}
		orderItems = append(orderItems, orderItem)
	}

	log.Printf("Create: Subtotal calculated: %.2f", subtotal)

	// Crear el pedido
	order := models.Order{
		UserID:      int(claim.User.ID),
		OrderNumber: generateOrderNumber(),
		Status:      "pending",
		TotalAmount: subtotal, // Por ahora sin shipping, tax, etc.
		Subtotal:    subtotal,
		Notes:       req.Notes,
	}

	if err := tx.Create(&order).Error; err != nil {
		log.Printf("Create: Error creating order: %v", err)
		tx.Rollback()
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Create: Order created successfully with ID: %d", order.ID)

	// Crear los order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := tx.Create(&orderItems[i]).Error; err != nil {
			log.Printf("Create: Error creating order item %d: %v", i, err)
			tx.Rollback()
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

	log.Printf("Create: Order items created successfully")

	// Crear información de pago
	payment := models.OrderPayment{
		OrderID:       order.ID,
		PaymentMethod: req.Payment.PaymentMethod,
		PaymentStatus: "pending",
		Amount:        req.Payment.Amount,
		Currency:      "ARS",
	}

	if err := tx.Create(&payment).Error; err != nil {
		log.Printf("Create: Error creating payment: %v", err)
		tx.Rollback()
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Create: Payment created successfully")

	// Crear información de envío
	shipping := models.OrderShipping{
		OrderID:            order.ID,
		RecipientName:      req.Shipping.RecipientName,
		RecipientEmail:     req.Shipping.RecipientEmail,
		RecipientPhone:     req.Shipping.RecipientPhone,
		ShippingAddress:    req.Shipping.ShippingAddress,
		ShippingCity:       req.Shipping.ShippingCity,
		ShippingPostalCode: req.Shipping.ShippingPostalCode,
		ShippingProvince:   req.Shipping.ShippingProvince,
		ShippingMethod:     req.Shipping.ShippingMethod,
		ShippingCountry:    "Argentina",
		ShippingStatus:     "pending",
	}

	if err := tx.Create(&shipping).Error; err != nil {
		log.Printf("Create: Error creating shipping: %v", err)
		tx.Rollback()
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Create: Shipping created successfully")

	// Commit de la transacción
	if err := tx.Commit().Error; err != nil {
		log.Printf("Create: Error committing transaction: %v", err)
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Create: Transaction committed successfully")

	// Preparar respuesta
	response := map[string]interface{}{
		"id":           order.ID,
		"order_number": order.OrderNumber,
		"status":       order.Status,
		"total_amount": order.TotalAmount,
		"created_at":   order.CreatedAt,
		"message":      "Pedido creado exitosamente",
	}

	log.Printf("Create: Order creation completed successfully")
	responses.JSON(w, http.StatusCreated, response)
}

func Index(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order

	// Obtener todos los pedidos con sus relaciones
	if err := db.Instance().Preload("User").Preload("OrderItems").Preload("Payment").Preload("Shipping").Find(&orders).Error; err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Convertir a respuesta
	var response []map[string]interface{}
	for _, order := range orders {
		orderData := map[string]interface{}{
			"id":           order.ID,
			"order_number": order.OrderNumber,
			"status":       order.Status,
			"total_amount": order.TotalAmount,
			"created_at":   order.CreatedAt,
			"user": map[string]interface{}{
				"id":    order.User.ID,
				"name":  order.User.Name,
				"email": order.User.Email,
			},
		}
		response = append(response, orderData)
	}

	responses.JSON(w, http.StatusOK, response)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var order models.Order
	if err := db.Instance().Preload("User").Preload("OrderItems").Preload("Payment").Preload("Shipping").First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.ERROR(w, http.StatusNotFound, fmt.Errorf("pedido no encontrado"))
			return
		}
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, order)
}

func GetByUser(w http.ResponseWriter, r *http.Request) {
	// Obtener el claim del contexto
	claim, ok := r.Context().Value(authentication.UserContextKey).(models.Claim)
	if !ok {
		responses.ERROR(w, http.StatusUnauthorized, nil)
		return
	}

	var orders []models.Order
	if err := db.Instance().Where("user_id = ?", claim.User.ID).Preload("OrderItems").Preload("Payment").Preload("Shipping").Find(&orders).Error; err != nil {
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

	// Validar status válido
	validStatuses := []string{"pending", "confirmed", "processing", "shipped", "delivered", "cancelled"}
	valid := false
	for _, status := range validStatuses {
		if status == req.Status {
			valid = true
			break
		}
	}
	if !valid {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("status inválido: %s", req.Status))
		return
	}

	// Actualizar el pedido
	var order models.Order
	if err := db.Instance().First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.ERROR(w, http.StatusNotFound, fmt.Errorf("pedido no encontrado"))
			return
		}
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	order.Status = req.Status
	now := time.Now()
	order.UpdatedAt = &now

	if err := db.Instance().Save(&order).Error; err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "Estado actualizado correctamente"})
}

func UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var req UpdatePaymentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Buscar el pago del pedido
	var payment models.OrderPayment
	if err := db.Instance().Where("order_id = ?", id).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.ERROR(w, http.StatusNotFound, fmt.Errorf("pago no encontrado para el pedido"))
			return
		}
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Actualizar estado y transaction_id si viene
	payment.PaymentStatus = req.PaymentStatus
	if req.TransactionID != "" {
		payment.TransactionID = req.TransactionID
	}

	if err := db.Instance().Save(&payment).Error; err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "Estado de pago actualizado correctamente"})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Verificar que el pedido existe
	var order models.Order
	if err := db.Instance().First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.ERROR(w, http.StatusNotFound, fmt.Errorf("pedido no encontrado"))
			return
		}
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Eliminar el pedido (soft delete)
	if err := db.Instance().Delete(&order).Error; err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "Pedido eliminado correctamente"})
}
