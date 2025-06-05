package models

import (
	"api-rbac/db"

	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	Products    []Product `json:"products"`
}

func (c *Category) FindAll() ([]Category, error) {
	db := db.Instance()

	categories := []Category{}

	err := db.Preload("Products").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
func (c *Category) GetByID(id int) (Category, error) {
	db := db.Instance()

	category := Category{}
	err := db.Preload("Products").Where("id = ?", id).First(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}
func (c *Category) Create() (Category, error) {
	db := db.Instance()

	err := db.Create(&c).Error
	if err != nil {
		return Category{}, err
	}

	return *c, nil
}
func (c *Category) Update() (Category, error) {
	db := db.Instance()

	err := db.Save(&c).Error
	if err != nil {
		return Category{}, err
	}

	return *c, nil
}
func (c *Category) Delete(id int) error {
	db := db.Instance()

	err := db.Where("id = ?", id).Delete(&c).Error
	if err != nil {
		return err
	}

	return nil
}
