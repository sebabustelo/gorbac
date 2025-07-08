package models

import (
	"api-rbac/db"

	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	UserID      uint        `json:"user_id"`
	User        User        `json:"user"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status" gorm:"default:'pending'"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

func (o *Order) Create() (Order, error) {
	db := db.Instance()
	err := db.Create(o).Error
	return *o, err
}

func (o *Order) FindAll() ([]Order, error) {
	db := db.Instance()
	orders := []Order{}
	err := db.Preload("User").Preload("Items.Product").Find(&orders).Error
	return orders, err
}

func (o *Order) GetByID(id int) (Order, error) {
	db := db.Instance()
	order := Order{}
	err := db.Preload("User").Preload("Items.Product").Where("id = ?", id).First(&order).Error
	return order, err
}

func (o *Order) GetByUser(userID int) ([]Order, error) {
	db := db.Instance()
	orders := []Order{}
	err := db.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (o *Order) UpdateStatus(id int, status string) error {
	db := db.Instance()
	return db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}

func (o *Order) Delete(id int) error {
	db := db.Instance()
	return db.Delete(&Order{}, id).Error
}
