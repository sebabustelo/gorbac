package models

import (
	"api-rbac/db"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User representa la estructura de usuario con sus atributos
type User struct {
	gorm.Model
	User       string    `json:"user" gorm:"unique;not null;size:50"`
	Email      string    `json:"email" gorm:"unique;not null;size:100"`
	Name       string    `json:"name" gorm:"size:50;not null"`
	LastName   string    `json:"last_name" gorm:"size:50"`
	Password   string    `json:"password,omitempty" gorm:"size:100;not null"`
	Provider   string    `json:"provider" gorm:"size:20;default:'local'"`
	ProviderID string    `json:"provider_id" gorm:"size:100"`
	LastLogin  time.Time `json:"last_login"`
	Active     bool      `json:"active" gorm:"default:true"`
	Roles      []Role    `json:"roles" gorm:"many2many:user_roles;"`
	Token      string    `json:"token,omitempty" gorm:"-" sql:"-"`
}

// TableName especifica el nombre de la tabla para el modelo User
func (User) TableName() string {
	return "users"
}

// BeforeSave hook de GORM para hashear la contraseña antes de guardar
func (u *User) BeforeSave() error {
	if u.Password != "" && len(u.Password) < 60 { // Solo hashear si no está ya hasheada
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error al hashear la contraseña: %v", err)
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword verifica si la contraseña proporcionada coincide
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// FindAll retorna todos los usuarios
func (u User) FindAll() ([]User, error) {
	db := db.Instance()
	users := []User{}
	err := db.Preload("Roles").Find(&users).Error
	return users, err
}

// Create crea un nuevo usuario
func (u *User) Create(user User) (*User, error) {
	db := db.Instance()

	// Verificar si el email o usuario ya existe
	var count int
	if err := db.Model(&User{}).Where("email = ? OR user = ?", user.Email, user.User).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("error al verificar duplicados: %v", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("el email o usuario ya existe")
	}

	// Crear el usuario
	if err := db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error al crear usuario: %v", err)
	}

	// Cargar las relaciones
	if err := db.Preload("Roles").First(&user, user.ID).Error; err != nil {
		return nil, fmt.Errorf("error al cargar relaciones: %v", err)
	}

	// Limpiar contraseña antes de retornar
	user.Password = ""
	return &user, nil
}

// Update actualiza un usuario existente
func (u *User) Update(user User) (*User, error) {
	db := db.Instance()

	// Verificar si el email o usuario ya existe en otro registro
	var count int
	if err := db.Model(&User{}).Where("(email = ? OR user = ?) AND id != ?", user.Email, user.User, user.ID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("el email o usuario ya existe en otro registro")
	}

	// Actualizar solo los campos que han cambiado
	updates := map[string]interface{}{
		"name":      user.Name,
		"last_name": user.LastName,
		"email":     user.Email,
		"user":      user.User,
		"active":    user.Active,
	}

	// Solo actualizar contraseña si se proporciona una nueva
	if user.Password != "" {
		updates["password"] = user.Password
	}

	err := db.Model(&User{}).Where("id = ?", user.ID).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Limpiar contraseña antes de retornar
	user.Password = ""
	return &user, nil
}

// Delete elimina un usuario
func (u *User) Delete(id int) (bool, error) {
	db := db.Instance()

	// Verificar si el usuario existe
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return false, fmt.Errorf("usuario no encontrado")
	}

	// Eliminar relaciones primero
	if err := db.Model(&user).Association("Roles").Clear().Error; err != nil {
		return false, err
	}

	// Eliminar usuario
	if err := db.Delete(&user).Error; err != nil {
		return false, err
	}

	return true, nil
}

// GetUserInfo recibe el nombre de usario y devuelve la informacion del mismo en la estructura User
func (u *User) GetUserInfo(usuario string) (User, error) {
	db := db.Instance()

	fmt.Println("GetUserInfo: ", usuario)
	user := User{}
	err := db.
		Where("user = ? OR email = ?", usuario, usuario).
		Preload("Roles").
		First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetByID ...
func (u *User) GetByID(id int) (User, error) {
	db := db.Instance()

	user := User{}
	err := db.Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// FindOrCreateSocial busca un usuario por email y provider, y lo crea si no existe.
func (u *User) FindOrCreateSocial(email, name, lastName, provider, providerID string) (*User, error) {
	db := db.Instance()
	var user User

	// Busca por email y provider
	err := db.Where("email = ? AND provider = ?", email, provider).First(&user).Error
	if err == nil {
		return &user, nil // Ya existe
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err // Otro error
	}

	// No existe, lo crea
	user = User{
		User:       email, // O usa otro campo único si prefieres
		Email:      email,
		Name:       name,
		LastName:   lastName,
		Provider:   provider,
		ProviderID: providerID,
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAllWithPagination retorna usuarios paginados y filtrados usando GORM
func (u User) FindAllWithPagination(page, limit int, search string) ([]User, int64, error) {
	db := db.Instance()
	var users []User
	var total int64

	// Crear el scope base
	scope := db.Model(&User{})

	// Aplicar búsqueda si se proporciona
	if search != "" {
		search = "%" + search + "%"
		scope = scope.Where(
			db.Where("name LIKE ?", search).
				Or("email LIKE ?", search).
				Or("user LIKE ?", search),
		)
	}

	// Obtener el total de registros
	if err := scope.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación y obtener resultados
	err := scope.
		Preload("Roles").
		Scopes(
			Paginate(page, limit),
			OrderByCreatedAt(),
		).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Paginate es un scope de GORM para manejar la paginación
func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

// OrderByCreatedAt es un scope de GORM para ordenar por fecha de creación
func OrderByCreatedAt() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}
}
