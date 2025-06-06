package models

import (
	"api-rbac/db"
	"testing"

	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestProduct_FindAll(t *testing.T) {
	SetupTestDB()

	product := Product{}
	products, err := product.FindAll()
	if err != nil {
		t.Errorf("Error finding all products: %v", err)
		return
	}

	if len(products) == 0 {
		t.Error("Expected to find products, but got none")
	}
}

func TestProduct_GetByID(t *testing.T) {
	SetupTestDB()

	product := Product{}
	// Assuming you have a product with ID 1 in your test database
	foundProduct, err := product.GetByID(1)
	if err != nil {
		t.Errorf("Error getting product by ID: %v", err)
		return
	}

	if foundProduct.ID != 1 {
		t.Errorf("Expected product ID 1, got %d", foundProduct.ID)
	}
}

func TestProduct_Create(t *testing.T) {
	SetupTestDB()

	// Obtener la categoría de prueba existente
	db := db.Instance()
	var category Category
	if err := db.Where("name = ?", "Test Category").First(&category).Error; err != nil {
		t.Fatalf("No se pudo obtener la categoría de prueba: %v", err)
	}

	product := &Product{
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       10.99,
		Stock:       100,
		CategoryID:  category.ID,
	}

	createdProduct, err := product.Create()
	if err != nil {
		t.Errorf("Error creating product: %v", err)
		return
	}

	if createdProduct.Name != product.Name {
		t.Errorf("Expected product name %s, got %s", product.Name, createdProduct.Name)
	}
}

func TestProduct_Update(t *testing.T) {
	SetupTestDB()

	category := &Category{
		Name:        "Updated Category",
		Description: "Categoría de prueba",
	}
	db := db.Instance()
	if err := db.Create(category).Error; err != nil {
		t.Fatalf("No se pudo crear la categoría de prueba: %v", err)
	}

	product := &Product{
		Model:       gorm.Model{ID: 1},
		Name:        "Updated Product",
		Description: "This is an updated product",
		Price:       12.99,
		Stock:       150,
		CategoryID:  category.ID,
	}

	updatedProduct, err := product.Update()
	if err != nil {
		t.Errorf("Error updating product: %v", err)
		return
	}

	if updatedProduct.Name != product.Name {
		t.Errorf("Expected updated product name %s, got %s", product.Name, updatedProduct.Name)
	}
}

// func TestProduct_Delete(t *testing.T) {

// 	product := &Product{
// 		Model: gorm.Model{ID: 1}, // Assuming product with ID 1 exists
// 	}

// 	err := product.Delete()
// 	if err != nil {
// 		t.Errorf("Error deleting product: %v", err)
// 		return
// 	}

// 	// Verify that the product was deleted
// 	_, err = product.GetByID(1)
// 	if err == nil {
// 		t.Error("Expected error when getting deleted product, but got none")
// 	}
// }
