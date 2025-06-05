package models

import (
	"api-rbac/db"

	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category"`
	Image       string   `json:"image"`
}

func (p *Product) FindAll() ([]Product, error) {
	db := db.Instance()

	products := []Product{}

	err := db.Preload("Category").Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (p *Product) GetByID(id int) (Product, error) {
	db := db.Instance()

	product := Product{}
	err := db.Preload("Category").Where("id = ?", id).First(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}
func (p *Product) Create() (Product, error) {
	db := db.Instance()

	err := db.Create(&p).Error
	if err != nil {
		return Product{}, err
	}

	return *p, nil
}
func (p *Product) Update() (Product, error) {
	db := db.Instance()

	err := db.Save(&p).Error
	if err != nil {
		return Product{}, err
	}

	return *p, nil
}
func (p *Product) Delete() error {
	db := db.Instance()

	err := db.Delete(&p).Error
	if err != nil {
		return err
	}

	return nil
}
func (p *Product) FindByCategoryID(categoryID uint) ([]Product, error) {
	db := db.Instance()

	products := []Product{}

	err := db.Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (p *Product) FindByName(name string) ([]Product, error) {
	db := db.Instance()

	products := []Product{}

	err := db.Where("name LIKE ?", "%"+name+"%").Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
