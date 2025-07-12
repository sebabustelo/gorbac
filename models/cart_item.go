package models

import (
	"api-rbac/db"

	"github.com/jinzhu/gorm"
)

type CartItem struct {
	gorm.Model
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}

func (ci *CartItem) GetByID(id int) (CartItem, error) {
	db := db.Instance()

	item := CartItem{}
	err := db.Preload("Product").Where("id = ?", id).First(&item).Error
	if err != nil {
		return item, err
	}

	return item, nil
}

func (ci *CartItem) GetByCartAndProduct(cartID, productID uint) (CartItem, error) {
	db := db.Instance()

	item := CartItem{}
	err := db.Preload("Product").Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		return item, err
	}

	return item, nil
}

func (ci *CartItem) GetByCartID(cartID uint) ([]CartItem, error) {
	db := db.Instance()

	items := []CartItem{}
	err := db.Preload("Product").Where("cart_id = ?", cartID).Find(&items).Error
	if err != nil {
		return items, err
	}

	return items, nil
}

func (ci *CartItem) Create() (CartItem, error) {
	db := db.Instance()

	err := db.Create(&ci).Error
	if err != nil {
		return CartItem{}, err
	}

	return *ci, nil
}

func (ci *CartItem) Update() (CartItem, error) {
	db := db.Instance()

	err := db.Save(&ci).Error
	if err != nil {
		return CartItem{}, err
	}

	return *ci, nil
}

func (ci *CartItem) Delete(id int) error {
	db := db.Instance()

	err := db.Where("id = ?", id).Delete(&ci).Error
	if err != nil {
		return err
	}

	return nil
}

func (ci *CartItem) DeleteByCartID(cartID uint) error {
	db := db.Instance()

	err := db.Where("cart_id = ?", cartID).Delete(&ci).Error
	if err != nil {
		return err
	}

	return nil
}
