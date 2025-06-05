package models

import (
	"api-rbac/db"
)

// SetupTestDB prepares the test database with initial data
func SetupTestDB() {
	db := db.TestInstance()

	// Verificar que estamos usando la base de datos correcta
	var dbName string
	db.Raw("SELECT DATABASE()").Row().Scan(&dbName)
	if dbName != "gorbac_test" {
		panic("No estamos usando la base de datos de prueba. Base de datos actual: " + dbName)
	}

	// Desactivar claves foráneas
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Limpiar tablas en orden correcto para evitar problemas de claves foráneas
	tables := []string{
		"cart_items",
		"carts",
		"products",
		"categories",
		"roles_apis",
		"user_roles",
		"apis",
		"roles",
		"users",
	}

	// Primero intentamos TRUNCATE
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}

	// Si TRUNCATE falla, intentamos DELETE
	for _, table := range tables {
		db.Exec("DELETE FROM " + table)
	}

	// Resetear auto-incrementos
	for _, table := range tables {
		db.Exec("ALTER TABLE " + table + " AUTO_INCREMENT = 1")
	}

	// Reactivar claves foráneas
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// Verificar que las tablas están vacías
	var count int64
	for _, table := range tables {
		db.Table(table).Count(&count)
		if count > 0 {
			panic("La tabla " + table + " no está vacía después de la limpieza")
		}
	}

	// Crear categoría de prueba
	category := &Category{
		Name:        "Test Category",
		Description: "Test Category Description",
	}
	db.Create(category)

	// Crear productos de prueba
	products := []Product{
		{
			Name:        "Test Product 1",
			Description: "Test Product 1 Description",
			Price:       10.99,
			Stock:       100,
			CategoryID:  category.ID,
		},
		{
			Name:        "Test Product 2",
			Description: "Test Product 2 Description",
			Price:       20.99,
			Stock:       200,
			CategoryID:  category.ID,
		},
	}
	for _, p := range products {
		db.Create(&p)
	}
}
