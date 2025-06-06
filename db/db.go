package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	"api-rbac/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db         *gorm.DB
	dbTest     *gorm.DB
	mu         sync.Mutex
	configPath string
)

// Instance returns the main database instance
func Instance() *gorm.DB {
	if db != nil {
		return db
	}

	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		return db
	}

	return connectDB()
}

// connectDB establishes a connection to the database
func connectDB() *gorm.DB {
	if os.Getenv("GO_ENV") == "test" {
		configPath = "./config/config_test.json"
	} else if os.Getenv("GO_ENV") == "local" {
		configPath = "./config/config_local.json"
	} else {
		configPath = "./config/config.json"
	}

	c, err := config.New(configPath)
	if err != nil {
		log.Fatalf("No se puede leer el archivo config.json en %s", configPath)
	}

	err = c.Validate("db_driver", "db_host", "db_port", "db_name", "db_user", "db_password")
	if err != nil {
		log.Fatal("No se puede validar los campos db_driver|db_host|db_port|db_name|db_user|db_password del archivo config.json")
	}

	dbHost, _ := c.Get("db_host")
	dbUser, _ := c.Get("db_user")
	dbPass, _ := c.Get("db_password")
	dbName, _ := c.Get("db_name")
	dbPort, _ := c.Get("db_port")

	// dbUser := "root"
	// dbPass := "wqwRdAcPeBlwQXWALkGPMIAzxXLclyAs"
	// dbHost := "mysql.railway.internal"
	// dbPort := 3306
	// dbName := "railway"

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	var connection *gorm.DB
	connection, err = gorm.Open("mysql", dbURI)
	if err != nil {
		log.Fatal(err)
	}

	db = connection

	return db
}
