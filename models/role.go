package models

import (
	"api-rbac/db"

	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	Name  string `json:"name" gorm:"unique;not null"`
	Apis  []Api  `gorm:"many2many:roles_apis;"`
	Users []User `gorm:"many2many:user_roles;"`
}

func (r *Role) FindAll() ([]Role, error) {
	db := db.Instance()

	var roles []Role
	// Preload las relaciones correctas
	err := db.Preload("Apis").Preload("Users").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Role) GetByID(id int) (Role, error) {
	db := db.Instance()

	var role Role
	err := db.Preload("Apis").Preload("Users").Where("id = ?", id).First(&role).Error
	if err != nil {
		return role, err
	}

	return role, nil
}

func (r *Role) GetApisByRole(roleID int) ([]Api, error) {
	var role Role
	err := db.Instance().Preload("Apis").Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return nil, err
	}
	return role.Apis, nil
}

func (r *Role) Create() (Role, error) {
	db := db.Instance()

	err := db.Create(&r).Error
	if err != nil {
		return Role{}, err
	}

	return *r, nil
}

func (r *Role) Update() (Role, error) {
	db := db.Instance()

	err := db.Save(&r).Error
	if err != nil {
		return Role{}, err
	}

	return *r, nil
}

func (r *Role) Delete() error {
	db := db.Instance()

	err := db.Delete(&r).Error
	if err != nil {
		return err
	}

	return nil
}
