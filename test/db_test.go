package test

import (
	"api-rbac/db"
	"testing"
)

func TestDatabaseConnections(t *testing.T) {
	// Probar conexión a base de datos principal
	mainDB := db.Instance()
	if mainDB == nil {
		t.Error("No se pudo conectar a la base de datos principal")
	}

	// Probar conexión a base de datos de prueba
	testDB := db.TestInstance()
	if testDB == nil {
		t.Error("No se pudo conectar a la base de datos de prueba")
	}

	// Verificar que son instancias diferentes
	if mainDB == testDB {
		t.Error("Las instancias de base de datos son iguales, deberían ser diferentes")
	}
}
