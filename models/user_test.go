package models

import (
	"api-rbac/db"
	"testing"

	// _ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
)

// TestUser_Create prueba la creación de usuarios
func TestUser_Create(t *testing.T) {
	SetupTestDB()

	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "crear usuario válido",
			user: User{
				User:     "testuser",
				Email:    "test@example.com",
				Name:     "Test",
				LastName: "User",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "crear usuario con email duplicado",
			user: User{
				User:     "testuser2",
				Email:    "test@example.com", // mismo email que el anterior
				Name:     "Test2",
				LastName: "User2",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "crear usuario con nombre de usuario duplicado",
			user: User{
				User:     "testuser", // mismo user que el primero
				Email:    "test2@example.com",
				Name:     "Test2",
				LastName: "User2",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{}
			got, err := u.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("User.Create() returned nil user when no error expected")
			}
		})
	}
}

// TestUser_Update prueba la actualización de usuarios
func TestUser_Update(t *testing.T) {
	SetupTestDB()

	// Crear usuario base para las pruebas
	baseUser := User{
		User:     "updateuser",
		Email:    "update@example.com",
		Name:     "Update",
		LastName: "User",
		Password: "password123",
	}
	u := &User{}
	created, err := u.Create(baseUser)
	if err != nil {
		t.Fatalf("Error al crear usuario base: %v", err)
	}

	// Crear otro usuario con el email que usaremos en la prueba
	otherUser := User{
		User:     "otheruser",
		Email:    "other@example.com", // Cambiado para evitar duplicados
		Name:     "Other",
		LastName: "User",
		Password: "password123",
	}
	_, err = u.Create(otherUser)
	if err != nil {
		t.Fatalf("Error al crear usuario adicional: %v", err)
	}

	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "actualizar nombre",
			user: User{
				Model:    gorm.Model{ID: created.Model.ID},
				User:     "updateuser",
				Email:    "update@example.com",
				Name:     "Updated",
				LastName: "User",
			},
			wantErr: false,
		},
		{
			name: "actualizar con email existente de otro usuario",
			user: User{
				Model:    gorm.Model{ID: created.Model.ID},
				User:     "updateuser",
				Email:    "other@example.com", // Usar el email del otro usuario
				Name:     "Update",
				LastName: "User",
			},
			wantErr: true,
		},
		{
			name: "actualizar contraseña",
			user: User{
				Model:    gorm.Model{ID: created.Model.ID},
				User:     "updateuser",
				Email:    "update@example.com",
				Name:     "Update",
				LastName: "User",
				Password: "newpassword123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.Update(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("User.Update() returned nil user when no error expected")
			}
		})
	}
}

// TestUser_Delete prueba la eliminación de usuarios
func TestUser_Delete(t *testing.T) {
	SetupTestDB()

	// Crear usuario para eliminar
	user := User{
		User:     "deleteuser",
		Email:    "delete@example.com",
		Name:     "Delete",
		LastName: "User",
		Password: "password123",
	}
	u := &User{}
	created, err := u.Create(user)
	if err != nil {
		t.Fatalf("Error al crear usuario: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "eliminar usuario existente",
			id:      int(created.ID),
			wantErr: false,
		},
		{
			name:    "eliminar usuario inexistente",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got {
				t.Error("User.Delete() returned false when no error expected")
			}
		})
	}
}

// TestUser_CheckPassword prueba la verificación de contraseñas
func TestUser_CheckPassword(t *testing.T) {
	SetupTestDB()

	// Crear usuario con contraseña
	user := User{
		User:     "passworduser",
		Email:    "password@example.com",
		Name:     "Password",
		LastName: "User",
		Password: "correctpassword",
	}
	u := &User{}
	created, err := u.Create(user)
	if err != nil {
		t.Fatalf("Error al crear usuario: %v", err)
	}

	// Verificar que el usuario se creó correctamente
	if created == nil {
		t.Fatal("Usuario creado es nil")
	}

	// Recargar el usuario para asegurar que tenemos la contraseña hasheada
	db := db.Instance()
	var reloadedUser User
	err = db.First(&reloadedUser, created.ID).Error
	if err != nil {
		t.Fatalf("Error al recargar usuario: %v", err)
	}

	// Verificar que el usuario recargado tiene los datos correctos
	if reloadedUser.ID != created.ID {
		t.Fatalf("ID del usuario recargado (%d) no coincide con el creado (%d)", reloadedUser.ID, created.ID)
	}

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "contraseña correcta",
			password: "correctpassword",
			want:     true,
		},
		{
			name:     "contraseña incorrecta",
			password: "wrongpassword",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reloadedUser.CheckPassword(tt.password); got != tt.want {
				t.Errorf("User.CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
