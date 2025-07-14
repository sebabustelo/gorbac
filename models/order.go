package models

import (
	"time"
)

// Order representa un pedido
type Order struct {
	ID             int        `json:"id" gorm:"primaryKey"`
	UserID         int        `json:"user_id" gorm:"not null"`
	OrderNumber    string     `json:"order_number" gorm:"unique;not null"`
	Status         string     `json:"status" gorm:"default:'pending'"`
	TotalAmount    float64    `json:"total_amount" gorm:"not null"`
	Subtotal       float64    `json:"subtotal" gorm:"not null"`
	ShippingCost   float64    `json:"shipping_cost" gorm:"default:0"`
	TaxAmount      float64    `json:"tax_amount" gorm:"default:0"`
	DiscountAmount float64    `json:"discount_amount" gorm:"default:0"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"index"`

	// Relaciones
	User       User          `json:"user" gorm:"foreignKey:UserID"`
	OrderItems []OrderItem   `json:"order_items" gorm:"foreignKey:OrderID"`
	Payment    OrderPayment  `json:"payment" gorm:"foreignKey:OrderID"`
	Shipping   OrderShipping `json:"shipping" gorm:"foreignKey:OrderID"`
}

// TableName especifica el nombre de la tabla para el modelo Order
func (Order) TableName() string {
	return "orders"
}

// OrderItem representa un ítem del pedido
type OrderItem struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	OrderID      int        `json:"order_id" gorm:"not null"`
	ProductID    int        `json:"product_id" gorm:"not null"`
	ProductName  string     `json:"product_name" gorm:"not null"`
	ProductImage string     `json:"product_image"`
	Quantity     int        `json:"quantity" gorm:"default:1"`
	UnitPrice    float64    `json:"unit_price" gorm:"not null"`
	TotalPrice   float64    `json:"total_price" gorm:"not null"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`

	// Relaciones
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// TableName especifica el nombre de la tabla para el modelo OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderPayment representa la información de pago
type OrderPayment struct {
	ID             int        `json:"id" gorm:"primaryKey"`
	OrderID        int        `json:"order_id" gorm:"not null"`
	PaymentMethod  string     `json:"payment_method" gorm:"not null"`
	PaymentStatus  string     `json:"payment_status" gorm:"default:'pending'"`
	TransactionID  string     `json:"transaction_id"`
	Amount         float64    `json:"amount" gorm:"not null"`
	Currency       string     `json:"currency" gorm:"default:'ARS'"`
	PaymentDate    *time.Time `json:"payment_date"`
	PaymentDetails string     `json:"payment_details" gorm:"type:json"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"index"`

	// Relaciones
	// Order Order `json:"order" gorm:"foreignKey:OrderID"`
}

// TableName especifica el nombre de la tabla para el modelo OrderPayment
func (OrderPayment) TableName() string {
	return "order_payments"
}

// OrderShipping representa la información de envío
type OrderShipping struct {
	ID                 int        `json:"id" gorm:"primaryKey"`
	OrderID            int        `json:"order_id" gorm:"not null"`
	RecipientName      string     `json:"recipient_name" gorm:"not null"`
	RecipientEmail     string     `json:"recipient_email" gorm:"not null"`
	RecipientPhone     string     `json:"recipient_phone"`
	ShippingAddress    string     `json:"shipping_address" gorm:"not null"`
	ShippingCity       string     `json:"shipping_city" gorm:"not null"`
	ShippingPostalCode string     `json:"shipping_postal_code"`
	ShippingProvince   string     `json:"shipping_province"`
	ShippingCountry    string     `json:"shipping_country" gorm:"default:'Argentina'"`
	ShippingMethod     string     `json:"shipping_method" gorm:"default:'standard'"`
	TrackingNumber     string     `json:"tracking_number"`
	EstimatedDelivery  *time.Time `json:"estimated_delivery"`
	ActualDeliveryDate *time.Time `json:"actual_delivery_date"`
	ShippingStatus     string     `json:"shipping_status" gorm:"default:'pending'"`
	ShippingNotes      string     `json:"shipping_notes"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" gorm:"index"`

	// Relaciones
	// Order Order `json:"order" gorm:"foreignKey:OrderID"`
}

// TableName especifica el nombre de la tabla para el modelo OrderShipping
func (OrderShipping) TableName() string {
	return "order_shipping"
}

// OrderRequest representa la solicitud de creación de pedido
type OrderRequest struct {
	UserID   int                `json:"user_id"`
	Items    []OrderItemRequest `json:"items"`
	Payment  PaymentRequest     `json:"payment"`
	Shipping ShippingRequest    `json:"shipping"`
	Notes    string             `json:"notes"`
}

// OrderItemRequest representa un ítem en la solicitud de pedido
type OrderItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// PaymentRequest representa la información de pago en la solicitud
type PaymentRequest struct {
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
}

// ShippingRequest representa la información de envío en la solicitud
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

// OrderResponse representa la respuesta completa de un pedido
type OrderResponse struct {
	ID             int        `json:"id"`
	OrderNumber    string     `json:"order_number"`
	Status         string     `json:"status"`
	TotalAmount    float64    `json:"total_amount"`
	Subtotal       float64    `json:"subtotal"`
	ShippingCost   float64    `json:"shipping_cost"`
	TaxAmount      float64    `json:"tax_amount"`
	DiscountAmount float64    `json:"discount_amount"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`

	User       User          `json:"user"`
	OrderItems []OrderItem   `json:"order_items"`
	Payment    OrderPayment  `json:"payment"`
	Shipping   OrderShipping `json:"shipping"`
}
