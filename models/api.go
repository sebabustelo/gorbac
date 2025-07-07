package models

import (
	"api-rbac/db"

	"github.com/jinzhu/gorm"
)

type Api struct {
	gorm.Model
	Endpoint    string `json:"endpoint" gorm:"unique;not null"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
	Public      bool   `json:"public"`
	Tipo        string `json:"tipo" gorm:"not null;default:'GET'"`
	Roles       []Role `gorm:"many2many:roles_apis;"`
}

func (a *Api) FindAll() ([]Api, error) {
	db := db.Instance()

	apis := []Api{}

	err := db.Preload("Roles").Find(&apis).Error
	if err != nil {
		return nil, err
	}

	return apis, nil
}

func (a *Api) GetByID(id int) (Api, error) {
	db := db.Instance()

	api := Api{}
	err := db.Preload("Roles").Where("id = ?", id).First(&api).Error
	if err != nil {
		return api, err
	}

	return api, nil
}

func (a *Api) Create() (Api, error) {
	db := db.Instance()

	err := db.Create(&a).Error
	if err != nil {
		return Api{}, err
	}

	return *a, nil
}

func (a *Api) Update() (Api, error) {
	db := db.Instance()

	err := db.Save(&a).Error
	if err != nil {
		return Api{}, err
	}

	return *a, nil
}

func (a *Api) Delete() error {
	db := db.Instance()

	err := db.Delete(&a).Error
	if err != nil {
		return err
	}

	return nil
}
