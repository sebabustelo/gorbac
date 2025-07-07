package main

import (
	"api-rbac/authentication"
	"api-rbac/controllers/apis"
	"api-rbac/controllers/products"
	"api-rbac/controllers/roles"
	"api-rbac/controllers/users"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/cors"
)

func main() {

	//ExampleLDAPClientAuthenticate()
	r := chi.NewRouter()

	// Get allowed origins from environment or use defaults
	allowedOrigins := []string{
		"https://juvapets.netlify.app",
		"http://localhost:3000",
		"http://localhost:5173",
		"http://localhost:8080",
		"http://localhost:8229",
	}

	// Add custom origins from environment variable
	if customOrigins := os.Getenv("ALLOWED_ORIGINS"); customOrigins != "" {
		// Split by comma and add to allowed origins
		for _, origin := range strings.Split(customOrigins, ",") {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	cors := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"Origin",
			"X-Requested-With",
		},
		ExposedHeaders:   []string{"Content-Length"},
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
		r.Post("/products", products.Add)
		r.Put("/products/{id}", products.Edit)
		r.Delete("/products/{id}", products.Delete)
		r.Get("/apis", apis.Index)
		r.Post("/apis", apis.Add)
		r.Put("/apis/{id}", apis.Edit)
		r.Delete("/apis/{id}", apis.Delete)

	})
	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authentication.Login)
		r.Get("/roles/permissions/{id}/apis", roles.GetApisByRole)
		//r.Get("/validate", authentication.TokenValid)
		r.Get("/refresh", authentication.RefreshToken)
		r.Get("/products", products.Index)
		r.Get("/apis", apis.Index)
		r.Post("/google-login", authentication.GoogleLogin)

		// Health check endpoint
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})

		// Auth check endpoint
		r.Get("/auth/check", func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"No token provided"}`))
				return
			}

			// Remove "Bearer " prefix
			token = strings.TrimPrefix(token, "Bearer ")

			// Basic token validation
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"Invalid token format"}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"token_present"}`))
		})

	})

	// Use Railway's PORT if available, otherwise use default 8229
	port := os.Getenv("PORT")
	if port == "" {
		port = "8229"
	}

	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, r)

}
