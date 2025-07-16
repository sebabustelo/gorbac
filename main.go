package main

import (
	"api-rbac/authentication"
	"api-rbac/controllers/apis"
	"api-rbac/controllers/cart"
	"api-rbac/controllers/categories"
	"api-rbac/controllers/orders"
	"api-rbac/controllers/products"
	"api-rbac/controllers/roles"
	"api-rbac/controllers/users"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/cors"
)

var (
	capturedRoutes []string
	routesMutex    sync.RWMutex
)

// routeCaptureMiddleware captura las rutas automáticamente
func routeCaptureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := r.URL.Path
		routesMutex.Lock()
		// Agregar ruta si no existe
		found := false
		for _, existingRoute := range capturedRoutes {
			if existingRoute == route {
				found = true
				break
			}
		}
		if !found {
			capturedRoutes = append(capturedRoutes, route)
		}
		routesMutex.Unlock()
		next.ServeHTTP(w, r)
	})
}

// getRoutes retorna las rutas capturadas dinámicamente
func getRoutes() []string {
	routesMutex.RLock()
	defer routesMutex.RUnlock()
	return capturedRoutes
}

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
		routeCaptureMiddleware,
	)
	r.Use(cors.Handler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authentication.Verifytoken)
		r.Use(authentication.AuthorizeEndpoint)

		// ==============================================
		// Endpoints de APIs Management
		// ==============================================
		r.Delete("/apis/{id}", apis.Delete)       // Eliminar un endpoint API específico
		r.Get("/apis", apis.Index)                // Listar todos los endpoints API
		r.Get("/apis/{id}", apis.GetByID)         // Obtener detalles de un endpoint API específico
		r.Post("/apis/add", apis.Add)             // Agregar un nuevo endpoint API
		r.Put("/apis/{id}", apis.Edit)            // Editar un endpoint API existente
		r.Get("/apis/sync", apis.SyncApis)        // Sincronizar endpoints (lectura)
		r.Post("/apis/sync", apis.AddMissingApis) // Sincronizar endpoints (escritura)

		// ==============================================
		// Endpoints de Carrito de Compras
		// ==============================================
		r.Delete("/cart", cart.ClearCart)                 // Vaciar el carrito por completo
		r.Get("/cart", cart.GetCart)                      // Obtener el contenido del carrito
		r.Post("/cart/add", cart.AddToCart)               // Agregar item al carrito
		r.Delete("/cart/items/{id}", cart.RemoveFromCart) // Eliminar item específico del carrito
		r.Put("/cart/items/{id}", cart.UpdateCartItem)    // Actualizar cantidad de item en carrito

		// ==============================================
		// Endpoints de Órdenes/Pedidos
		// ==============================================
		r.Post("/orders", orders.Create)                          // Crear una nueva orden
		r.Get("/orders", orders.Index)                            // Listar todas las órdenes
		r.Get("/orders/{id}", orders.GetByID)                     // Obtener detalles de una orden específica
		r.Get("/orders/user", orders.GetByUser)                   // Obtener órdenes del usuario actual
		r.Delete("/orders/{id}", orders.Delete)                   // Eliminar una orden (admin)
		r.Put("/orders/{id}/payment", orders.UpdatePaymentStatus) // Actualizar estado de pago
		r.Put("/orders/{id}/status", orders.UpdateStatus)         // Actualizar estado de la orden

		// ==============================================
		// Endpoints de Productos
		// ==============================================
		r.Post("/products/add", products.Add)       // Agregar un nuevo producto
		r.Delete("/products/{id}", products.Delete) // Eliminar un producto
		r.Put("/products/{id}", products.Edit)      // Editar un producto existente

		// ==============================================
		// Endpoints de Roles
		// ==============================================
		r.Get("/roles", roles.Index)                // Listar todos los roles
		r.Post("/roles/add", roles.Add)             // Crear un nuevo rol
		r.Get("/roles/{id}", roles.GetByID)         // Obtener detalles de un rol específico
		r.Put("/roles/{id}/apis", roles.UpdateApis) // Actualizar APIs asociadas a un rol

		// ==============================================
		// Endpoints de Usuarios
		// ==============================================
		r.Post("/users/add", users.Add)              // Registrar un nuevo usuario
		r.Post("/users/edit", users.Edit)            // Editar información de usuario
		r.Delete("/users/delete/{id}", users.Delete) // Eliminar un usuario (admin)
		r.Get("/users", users.Index)                 // Listar todos los usuarios (admin)
		r.Get("/users/{id}", users.GetByID)          // Obtener información de usuario específico

	})
	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authentication.Login)
		r.Get("/roles/permissions/{id}/apis", roles.GetApisByRole)
		//r.Get("/validate", authentication.TokenValid)
		r.Get("/refresh", authentication.RefreshToken)
		r.Get("/products", products.Index)
		r.Get("/products/{id}", products.GetByID)
		r.Post("/google-login", authentication.GoogleLogin)
		// Public orders routes - requires authentication
		r.Get("/categories", categories.GetCategories)

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

	// Establecer las rutas dinámicamente para el sistema de APIs
	apis.SetRoutes(getRoutes())

	// Use Railway's PORT if available, otherwise use default 8229
	port := os.Getenv("PORT")
	if port == "" {
		port = "8229"
	}

	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, r)

}
