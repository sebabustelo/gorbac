package config_test

import (
	"api-rbac/config"
	"os"
	"path/filepath"
	"testing"
)

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = filepath.Join("config", "config_test.json")
	}
	return configPath
}

func TestLoadFile(t *testing.T) {
	configPath := getConfigPath()
	_, err := config.LoadFile(configPath)
	if err != nil {
		t.Fatalf("no se pudo leer el archivo %s: %v", configPath, err)
	}
}

func TestLoadBytes(t *testing.T) {
	j := []byte(`{
		"db_driver": "mysql",
		"db_host": "localhost",
		"db_port": "3306",
		"db_name": "test_db",
		"db_user": "test_user",
		"db_password": "test_pass"
	}`)

	c := config.Configuration{}
	err := config.LoadBytes(j, &c)
	if err != nil {
		t.Fatalf("no se pudo leer los bytes: %v", err)
	}

	// Verificar que los valores se cargaron correctamente
	driver, err := c.Get("db_driver")
	if err != nil {
		t.Errorf("no se pudo obtener db_driver: %v", err)
	}
	if driver != "mysql" {
		t.Errorf("db_driver = %v, want %v", driver, "mysql")
	}
}

func TestNew(t *testing.T) {
	configPath := getConfigPath()
	c, err := config.New(configPath)
	if err != nil {
		t.Fatalf("no se pudo cargar la configuracion: %v", err)
	}
	if c == nil {
		t.Fatal("la configuración no debería ser nil")
	}
}

func TestConfiguration_Validate(t *testing.T) {
	configPath := getConfigPath()
	c, err := config.New(configPath)
	if err != nil {
		t.Fatalf("no se pudo cargar la configuracion: %v", err)
	}

	requiredFields := []string{"db_driver", "db_host", "db_port", "db_name", "db_user", "db_password"}
	err = c.Validate(requiredFields...)
	if err != nil {
		t.Errorf("validación falló: %v", err)
	}
}

func TestConfiguration_Get(t *testing.T) {
	configPath := getConfigPath()
	c, err := config.New(configPath)
	if err != nil {
		t.Fatalf("no se pudo cargar la configuracion: %v", err)
	}

	value, err := c.Get("db_host")
	if err != nil {
		t.Errorf("no se pudo obtener db_host: %v", err)
	}
	if value == "" {
		t.Error("db_host no debería estar vacío")
	}
}

func TestConfiguration_GetInt(t *testing.T) {
	configPath := getConfigPath()
	c, err := config.New(configPath)
	if err != nil {
		t.Fatalf("no se pudo cargar la configuracion: %v", err)
	}

	port, err := c.GetInt("db_port")
	if err != nil {
		t.Errorf("no se pudo obtener db_port como entero: %v", err)
	}
	if port <= 0 {
		t.Errorf("db_port debería ser un número positivo, got %d", port)
	}
}
