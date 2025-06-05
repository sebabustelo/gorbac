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
	Roles       []Role `gorm:"many2many:roles_apis;"`
}

func (u *Api) FindAll() ([]Api, error) {
	db := db.Instance()

	apis := []Api{}

	db.Preload("Role").Find(&apis)

	return apis, nil
}
