package main

import (
	"api-rbac/authentication"
	"api-rbac/controllers/products"
	"api-rbac/controllers/roles"
	"api-rbac/controllers/users"
	"database/sql"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/cors"
)

func main() {

	//ExampleLDAPClientAuthenticate()
	r := chi.NewRouter()

	cors := cors.New(cors.Options{

		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Accept-Encoding", "Authorization", "Content-Type", "X-CSRF-Token"},

		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.Compress(5, "gzip"),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)
	r.Use(cors.Handler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authentication.Verifytoken)
		r.Use(authentication.AuthorizeEndpoint)

		r.Post("/users/add", users.Add)
		r.Post("/users/edit", users.Edit)
		r.Delete("/users/delete/{id}", users.Delete)
		r.Get("/users/{id}", users.GetByID)
		r.Get("/users/index", users.Index)
		r.Get("/roles", roles.Index)
		r.Post("/roles/add", roles.Add)
		r.Get("/roles/{id}", roles.GetByID)

	})
	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authentication.Login)
		r.Get("/roles/permissions/{id}/apis", roles.GetApisByRole)
		//r.Get("/validate", authentication.TokenValid)
		r.Get("/refresh", authentication.RefreshToken)
		r.Get("/products", products.Index)
		r.Post("/google-login", authentication.GoogleLogin)

	})

	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQLHOST")      // Railway suele usar "mysql.railway.internal"
	dbPort := os.Getenv("MYSQL_PORT")     // Generalmente "3306"
	dbName := os.Getenv("MYSQL_DATABASE") // normalmente "railway"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al abrir la conexión:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error al conectar a la base:", err)
	}

	fmt.Println("✅ Conexión exitosa a MySQL en Railway")

	http.ListenAndServe(":8229", r)

}
