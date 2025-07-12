package cart

import (
	"api-rbac/authentication"
	"api-rbac/models"
	"api-rbac/responses"
	"encoding/json"
	"net/http"
	"strconv"

	"fmt"

	"github.com/go-chi/chi"
)

type AddToCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

// GetCart obtiene el carrito del usuario
func GetCart(w http.ResponseWriter, r *http.Request) {
	claim, ok := r.Context().Value(authentication.UserContextKey).(models.Claim)
	if !ok {
		fmt.Println("❌ No se encontró el claim en el contexto")
		responses.ERROR(w, http.StatusUnauthorized, nil)
		return
	}
	fmt.Printf("✅ Claim encontrado: %+v\n", claim)

	userID := claim.User.ID

	cart := models.Cart{}
	err := cart.GetByUserID(userID)
	if err != nil {
		// Si no existe carrito, devolver carrito vacío
		emptyCart := models.Cart{
			UserID:    userID,
			CartItems: []models.CartItem{},
		}
		responses.JSON(w, http.StatusOK, emptyCart)
		return
	}

	responses.JSON(w, http.StatusOK, cart)
}

// AddToCart agrega un producto al carrito
func AddToCart(w http.ResponseWriter, r *http.Request) {
	// Obtener el claim del contexto
	claim, ok := r.Context().Value(authentication.UserContextKey).(models.Claim)
	if !ok {
		responses.ERROR(w, http.StatusUnauthorized, nil)
		return
	}

	userID := claim.User.ID

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if req.ProductID == 0 || req.Quantity <= 0 {
		responses.ERROR(w, http.StatusBadRequest, nil)
		return
	}

	// Verificar que el producto existe
	product := models.Product{}
	productFound, err := product.GetByID(int(req.ProductID))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	// Verificar stock
	if productFound.Stock < req.Quantity {
		responses.ERROR(w, http.StatusBadRequest, nil)
		return
	}

	// Obtener o crear carrito del usuario
	cart := models.Cart{}
	err = cart.GetByUserID(userID)
	if err != nil {
		// Crear nuevo carrito
		cart.UserID = userID
		cart, err = cart.Create()
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

	// Verificar si el producto ya está en el carrito
	cartItem := models.CartItem{}
	existingItem, err := cartItem.GetByCartAndProduct(cart.ID, req.ProductID)
	if err == nil {
		// Actualizar cantidad
		existingItem.Quantity += req.Quantity
		existingItem.Price = productFound.Price
		_, err = existingItem.Update()
	} else {
		// Crear nuevo item
		cartItem.CartID = cart.ID
		cartItem.ProductID = req.ProductID
		cartItem.Quantity = req.Quantity
		cartItem.Price = productFound.Price
		_, err = cartItem.Create()
	}

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Actualizar stock del producto
	productFound.Stock -= req.Quantity
	productFound.Update()

	// Obtener carrito actualizado
	updatedCart, _ := cart.GetByID(int(cart.ID))
	responses.JSON(w, http.StatusOK, updatedCart)
}

// UpdateCartItem actualiza la cantidad de un item en el carrito
func UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	// Obtener el claim del contexto
	claim, ok := r.Context().Value(authentication.UserContextKey).(models.Claim)
	if !ok {
		responses.ERROR(w, http.StatusUnauthorized, nil)
		return
	}

	userID := claim.User.ID
	itemID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var req UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if req.Quantity <= 0 {
		responses.ERROR(w, http.StatusBadRequest, nil)
		return
	}

	// Verificar que el item pertenece al usuario
	cartItem := models.CartItem{}
	item, err := cartItem.GetByID(itemID)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	cart := models.Cart{}
	cartFound, err := cart.GetByID(int(item.CartID))
	if err != nil || cartFound.UserID != userID {
		responses.ERROR(w, http.StatusForbidden, nil)
		return
	}

	// Verificar stock
	product := models.Product{}
	productFound, err := product.GetByID(int(item.ProductID))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	// Calcular diferencia de cantidad
	diff := req.Quantity - item.Quantity
	if productFound.Stock < diff {
		responses.ERROR(w, http.StatusBadRequest, nil)
		return
	}

	// Actualizar item
	item.Quantity = req.Quantity
	item.Price = productFound.Price
	_, err = item.Update()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Actualizar stock
	productFound.Stock -= diff
	productFound.Update()

	// Obtener carrito actualizado
	updatedCart, _ := cart.GetByID(int(cart.ID))
	responses.JSON(w, http.StatusOK, updatedCart)
}

// RemoveFromCart elimina un item del carrito
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	// Obtener el claim del contexto
	claim, ok := r.Context().Value(authentication.UserContextKey).(models.Claim)
	if !ok {
		responses.ERROR(w, http.StatusUnauthorized, nil)
		return
	}

	userID := claim.User.ID
	itemID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Verificar que el item pertenece al usuario
	cartItem := models.CartItem{}
	item, err := cartItem.GetByID(itemID)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	cart := models.Cart{}
	cartFound, err := cart.GetByID(int(item.CartID))
	if err != nil || cartFound.UserID != userID {
		responses.ERROR(w, http.StatusForbidden, nil)
		return
	}

	// Restaurar stock
	product := models.Product{}
	productFound, err := product.GetByID(int(item.ProductID))
	if err == nil {
		productFound.Stock += item.Quantity
		productFound.Update()
	}

	// Eliminar item
	err = cartItem.Delete(itemID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Obtener carrito actualizado
	updatedCart, _ := cart.GetByID(int(cart.ID))
	responses.JSON(w, http.StatusOK, updatedCart)
}

// ClearCart limpia todo el carrito del usuario
func ClearCart(w http.ResponseWriter, r *http.Request) {
	// Obtener el claim del contexto
	claim, ok := r.Context().Value(authentication.UserContextKey).(models.Claim)
	if !ok {
		responses.ERROR(w, http.StatusUnauthorized, nil)
		return
	}

	userID := claim.User.ID

	cart := models.Cart{}
	err := cart.GetByUserID(userID)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	// Restaurar stock de todos los items
	cartItems := models.CartItem{}
	items, err := cartItems.GetByCartID(cart.ID)
	if err == nil {
		for _, item := range items {
			product := models.Product{}
			productFound, err := product.GetByID(int(item.ProductID))
			if err == nil {
				productFound.Stock += item.Quantity
				productFound.Update()
			}
		}
	}

	// Eliminar todos los items
	err = cartItems.DeleteByCartID(cart.ID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "Carrito limpiado correctamente"})
}
