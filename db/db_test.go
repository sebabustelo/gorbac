package db

import (
	"os"
	"path/filepath"
	"testing"
)

func init() {
	// Configurar el archivo de configuración de prueba
	configPath := filepath.Join("/opt/go", "config", "config_test.json")
	os.Setenv("TEST_CONFIG_PATH", configPath)
}

func TestMainDBConnection(t *testing.T) {
	// Probar conexión a base de datos principal
	db := Instance()
	if db == nil {
		t.Error("No se pudo conectar a la base de datos principal")
	}
}

func TestTestDBConnection(t *testing.T) {
	// Probar conexión a base de datos de prueba
	db := Instance()
	if db == nil {
		t.Error("No se pudo conectar a la base de datos de prueba")
	}

	// Verificar que se conectó a la base de datos correcta
	var result string
	db.Raw("SELECT DATABASE()").Row().Scan(&result)
	if result != "gorbac_test" {
		t.Errorf("Se conectó a %s, se esperaba gorbac_test", result)
	}
}
