package models

import (
	"api-rbac/db"

	"fmt"

	"github.com/jinzhu/gorm"
)

type Cart struct {
	gorm.Model
	UserID    uint       `json:"user_id"`
	User      User       `json:"user"`
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID"`
}

func (c *Cart) FindAll() ([]Cart, error) {
	db := db.Instance()

	carts := []Cart{}

	err := db.Preload("User").Preload("CartItems.Product").Find(&carts).Error
	if err != nil {
		return nil, err
	}

	return carts, nil
}
func (c *Cart) GetByUserID(userID uint) error {
	db := db.Instance()

	fmt.Printf("🔍 Buscando carrito para userID: %d\n", userID)

	err := db.Where("user_id = ?", userID).
		Preload("User").
		Preload("CartItems.Product").
		First(c).Error

	if err != nil {
		fmt.Printf("❌ Error al buscar carrito: %v\n", err)
		return err
	}

	fmt.Printf("✅ Carrito encontrado - ID: %d, Items: %d\n", c.ID, len(c.CartItems))

	return nil
}

func (c *Cart) GetByID(id int) (Cart, error) {
	db := db.Instance()

	cart := Cart{}
	err := db.Preload("User").Preload("CartItems.Product").Where("id = ?", id).First(&cart).Error
	if err != nil {
		return cart, err
	}

	return cart, nil
}
func (c *Cart) Create() (Cart, error) {
	db := db.Instance()

	err := db.Create(&c).Error
	if err != nil {
		return Cart{}, err
	}

	return *c, nil
}
func (c *Cart) Update() (Cart, error) {
	db := db.Instance()

	err := db.Save(&c).Error
	if err != nil {
		return Cart{}, err
	}

	return *c, nil
}
func (c *Cart) Delete(id int) error {
	db := db.Instance()

	err := db.Where("id = ?", id).Delete(&c).Error
	if err != nil {
		return err
	}

	return nil
}
