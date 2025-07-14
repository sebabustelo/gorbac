package authentication

import (
	"api-rbac/db"
	"api-rbac/models"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	responses "api-rbac/helpers"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"golang.org/x/crypto/bcrypt"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	//privateBytes string
)

func init() {

	privateBytes, err := os.ReadFile("./private.rsa")
	if err != nil {
		log.Fatal("No se puede leer el arhivo privado")

	}

	publicBytes, err := os.ReadFile("./public.rsa.pub")
	if err != nil {
		log.Fatal("No se puede leer el arhivo publico")

	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se puedo hacer el parse a privatekey")

	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("No se puedo hacer el parse a publickey")

	}
}

// GenerateJWT token
func GenerateJWT(user models.User) string {
	claims := models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "RBAC",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatal("No se pudo firmar el token")
	}

	return result
}

// HashPassword ...
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash ...
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login verifica el user y password
func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("datos de login inválidos"))
		return
	}

	userInfo, err := userInput.GetUserInfo(userInput.User)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("usuario no encontrado"))
		return
	}

	// Login local
	if !CheckPasswordHash(userInput.Password, userInfo.Password) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("usuario y/o clave inválidas"))
		return
	}

	// Genera token y responde
	userInfo.Password = ""
	userInfo.Token = GenerateJWT(userInfo)
	responses.JSON(w, http.StatusOK, userInfo)
}

// Login con Google
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IdToken string `json:"id_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decodificando body:", err)
		responses.ERROR(w, http.StatusBadRequest, errors.New("datos de login inválidos"))
		return
	}

	if req.IdToken == "" {
		log.Println("GoogleLogin: id_token vacío")
		responses.ERROR(w, http.StatusBadRequest, errors.New("id_token requerido"))
		return
	}

	log.Printf("GoogleLogin: id_token recibido (longitud: %d)", len(req.IdToken))

	// Decodifica el token para ver el aud
	parts := strings.Split(req.IdToken, ".")
	if len(parts) == 3 {
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err == nil {
			var claims map[string]interface{}
			json.Unmarshal(payload, &claims)
			log.Printf("GoogleLogin: aud claim en el token: %v", claims["aud"])
			log.Printf("GoogleLogin: iss claim en el token: %v", claims["iss"])
			log.Printf("GoogleLogin: sub claim en el token: %v", claims["sub"])
			log.Printf("GoogleLogin: email claim en el token: %v", claims["email"])
		} else {
			log.Printf("GoogleLogin: Error decodificando payload: %v", err)
		}
	} else {
		log.Printf("GoogleLogin: Token no tiene el formato esperado (3 partes)")
	}

	// 1. Verifica el token de Google usando validación manual
	token, err := validateGoogleToken(req.IdToken)
	if err != nil {
		log.Printf("Error en validateGoogleToken: %v", err)
		responses.ERROR(w, http.StatusUnauthorized, errors.New("token de Google inválido"))
		return
	}
	log.Println("GoogleLogin: validateGoogleToken OK")

	// 2. Extrae datos del usuario de Google desde el token validado
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Error: claims inválidos")
		responses.ERROR(w, http.StatusUnauthorized, errors.New("claims del token inválidos"))
		return
	}

	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)
	sub, _ := claims["sub"].(string) // Google user ID
	log.Printf("GoogleLogin: email=%s, name=%s, sub=%s\n", email, name, sub)

	// 3. Busca o crea el usuario en tu base de datos
	var user models.User
	db := db.Instance()
	err = db.Where("email = ?", email).First(&user).Preload("Roles").Error
	if err != nil {
		log.Println("GoogleLogin: usuario no encontrado, creando uno nuevo")
		// Si no existe, créalo
		user = models.User{
			User:       email, // Usar email como user para evitar duplicados
			Email:      email,
			Name:       name,
			Provider:   "google",
			ProviderID: sub,
			Active:     true,
		}
		// Asigna el rol por defecto (rol ID 2 que de user)
		var defaultRole models.Role
		if err := db.First(&defaultRole, 2).Error; err == nil {
			user.Roles = []models.Role{defaultRole}
		}
		if err := db.Create(&user).Error; err != nil {
			log.Println("Error al crear usuario:", err)
			// Si es un error de duplicado, intentar buscar el usuario existente
			if strings.Contains(err.Error(), "Duplicate entry") {
				log.Println("Usuario duplicado detectado, buscando usuario existente...")
				var existingUser models.User
				if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
					user = existingUser
					log.Println("Usuario existente encontrado:", user.Email)
				} else {
					responses.ERROR(w, http.StatusInternalServerError, errors.New("error al buscar usuario existente"))
					return
				}
			} else {
				responses.ERROR(w, http.StatusInternalServerError, errors.New("error al crear el usuario"))
				return
			}
		}
	} else {
		log.Println("GoogleLogin: usuario encontrado en la base de datos")
		// Actualizar el nombre si está vacío o es diferente
		if user.Name == "" || user.Name != name {
			user.Name = name
			db.Model(&user).Update("name", name)
		}
		// Actualizar el campo User si está vacío
		if user.User == "" {
			user.User = email
			db.Model(&user).Update("user", email)
		}
		// Asegurar que tiene roles
		if len(user.Roles) == 0 {
			var defaultRole models.Role
			if err := db.First(&defaultRole, 2).Error; err == nil {
				db.Model(&user).Association("Roles").Append(defaultRole)
			}
		}
	}

	// 4. Genera tu JWT y responde
	user.Password = ""
	user.Token = GenerateJWT(user)
	log.Printf("GoogleLogin: Token generado: %s", user.Token)
	log.Printf("GoogleLogin: Usuario completo a enviar: %+v", user)
	log.Println("GoogleLogin: usuario autenticado, enviando respuesta OK")
	responses.JSON(w, http.StatusOK, user)
}

// RefreshToken Genera un nuevo token a partir de un token valido y lo entrega
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenValidation := TokenValid(r)

	if tokenValidation {
		tokenString := strings.Split(r.Header["Authorization"][0], "Bearer ")
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString[1], jwt.MapClaims{})
		if err != nil {
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
			token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
			result, err := token.SignedString(privateKey)
			if err != nil {
				log.Fatal("No se pudo firmar el token")
			}
			jsonResult, err := json.Marshal(result)
			if err != nil {
				log.Fatal("No se pudo firmar el token")
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "aplication/json")
			w.Write(jsonResult)
		} else {
			return
		}

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Su token no es válido")
		return
	}
}

// ValidateToken valida
func Verifytoken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details for debugging
		log.Printf("Verifytoken: Processing request to %s %s", r.Method, r.URL.Path)

		// Log Authorization header
		authHeader := r.Header.Get("Authorization")
		log.Printf("Verifytoken: Authorization header: %s", authHeader)

		token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			log.Printf("Verifytoken: Token parsing error: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Su token no es válido")
			return
		}

		if !token.Valid {
			log.Printf("Verifytoken: Token is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Su token no es válido")
			return
		}

		// Extrae el claim y lo guarda en el contexto
		if claims, ok := token.Claims.(*models.Claim); ok {
			log.Printf("Verifytoken: Token valid for user ID: %d", claims.User.ID)
			ctx := context.WithValue(r.Context(), UserContextKey, *claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		log.Printf("Verifytoken: Failed to extract claims from token")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Su token no es válido")
	})
}

// TokenValid valida el token
func TokenValid(r *http.Request) (status bool) {

	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		switch vErr := err.(type) {
		case *jwt.ValidationError:
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				return false
			case jwt.ValidationErrorClaimsInvalid:
				return false
			default:
				return false
			}
		default:
			return
		}
	}

	if token.Valid {
		status = true
	} else {
		status = false
	}

	return status

}

// Función para validar el token de Firebase/Google usando validación manual mejorada
func validateGoogleToken(idToken string) (*jwt.Token, error) {
	// Decodificar el token sin verificar para extraer claims
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Log todos los claims para debugging
	log.Printf("FirebaseLogin: Claims del token: %+v", claims)

	// Verificar audience (client ID de Firebase)
	aud, ok := claims["aud"].(string)
	if !ok {
		return nil, fmt.Errorf("missing audience claim")
	}

	// Firebase puede tener múltiples audience válidos
	validAudiences := []string{"petplace-dv69c", "petplace-dv69c.appspot.com"}
	validAudience := false
	for _, validAud := range validAudiences {
		if aud == validAud {
			validAudience = true
			break
		}
	}

	if !validAudience {
		return nil, fmt.Errorf("invalid audience: expected one of %v, got '%s'", validAudiences, aud)
	}

	// Verificar issuer (debe ser de Google/Firebase)
	iss, ok := claims["iss"].(string)
	if !ok {
		return nil, fmt.Errorf("missing issuer claim")
	}

	// Firebase puede usar diferentes issuers
	validIssuers := []string{
		"https://securetoken.google.com/petplace-dv69c",
		"https://accounts.google.com",
	}

	validIssuer := false
	for _, validIss := range validIssuers {
		if iss == validIss {
			validIssuer = true
			break
		}
	}

	if !validIssuer {
		return nil, fmt.Errorf("invalid issuer: expected one of %v, got '%s'", validIssuers, iss)
	}

	// Verificar expiración
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing expiration claim")
	}

	currentTime := float64(time.Now().Unix())
	if currentTime > exp {
		return nil, fmt.Errorf("token expired (exp: %f, current: %f)", exp, currentTime)
	}

	// Verificar que el token no haya sido emitido en el futuro (con tolerancia de 5 minutos)
	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing issued at claim")
	}

	tolerance := float64(5 * 60) // 5 minutos en segundos
	if currentTime < (iat - tolerance) {
		return nil, fmt.Errorf("token issued too far in the future (iat: %f, current: %f)", iat, currentTime)
	}

	// Verificar que tenga un subject (user ID)
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return nil, fmt.Errorf("missing or invalid subject claim")
	}

	// Verificar que tenga un email
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("missing or invalid email claim")
	}

	// Verificar que el email esté verificado (opcional para Firebase)
	emailVerified, ok := claims["email_verified"].(bool)
	if !ok {
		log.Printf("Warning: email_verified claim missing, assuming verified")
		emailVerified = true
	}

	if !emailVerified {
		log.Printf("Warning: email not verified for user: %s", email)
		// No bloqueamos por esto, solo log
	}

	// Verificar auth_time (tiempo de autenticación de Firebase)
	authTime, ok := claims["auth_time"].(float64)
	if !ok {
		log.Printf("Warning: auth_time claim missing")
	} else {
		// auth_time no debe ser muy antiguo (máximo 1 hora)
		maxAuthTime := currentTime - 3600 // 1 hora
		if authTime < maxAuthTime {
			return nil, fmt.Errorf("auth_time too old (auth_time: %f, max: %f)", authTime, maxAuthTime)
		}
	}

	// Verificar firebase claims específicos
	if firebaseClaims, ok := claims["firebase"].(map[string]interface{}); ok {
		log.Printf("FirebaseLogin: Firebase claims: %+v", firebaseClaims)

		// Verificar sign_in_provider
		if signInProvider, ok := firebaseClaims["sign_in_provider"].(string); ok {
			log.Printf("FirebaseLogin: Sign in provider: %s", signInProvider)
		}
	}

	log.Printf("FirebaseLogin: Token validation successful for user: %s (sub: %s)", email, sub)
	return token, nil
}

// Añadir función para comparar endpoints con parámetros
func matchEndpoint(pattern, path string) bool {
	re := regexp.MustCompile(`\{[^/]+\}`)
	patternRegex := re.ReplaceAllString(pattern, `[^/]+`)
	matched, _ := regexp.MatchString("^"+patternRegex+"$", path)
	return matched
}

// AuthorizeEndpoint verifica si el usuario tiene permiso para acceder al endpoint
func AuthorizeEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AuthorizeEndpoint: Processing request to %s %s", r.Method, r.URL.Path)

		claim, ok := r.Context().Value(UserContextKey).(models.Claim)
		if !ok {
			log.Printf("AuthorizeEndpoint: No claim found in context")
			http.Error(w, "no autorizado", http.StatusUnauthorized)
			return
		}
		userID := claim.User.ID
		endpoint := r.URL.Path
		method := r.Method

		log.Printf("AuthorizeEndpoint: User ID: %d, Endpoint: %s, Method: %s", userID, endpoint, method)

		db := db.Instance()

		var user models.User
		err := db.Where("id = ?", userID).
			Preload("Roles").
			Preload("Roles.Apis").
			First(&user).Error
		if err != nil {
			log.Printf("AuthorizeEndpoint: User not found: %v", err)
			http.Error(w, "usuario no encontrado", http.StatusUnauthorized)
			return
		}

		log.Printf("AuthorizeEndpoint: User has %d roles", len(user.Roles))

		hasPermission := false
		for _, role := range user.Roles {
			log.Printf("AuthorizeEndpoint: Checking role: %s (ID: %d)", role.Name, role.ID)
			for _, api := range role.Apis {
				log.Printf("AuthorizeEndpoint: Checking API: %s %s", api.Tipo, api.Endpoint)
				if matchEndpoint(api.Endpoint, endpoint) && strings.EqualFold(api.Tipo, method) {
					log.Printf("AuthorizeEndpoint: Permission granted for %s %s", method, endpoint)
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			log.Printf("AuthorizeEndpoint: Permission denied for %s %s", method, endpoint)
			log.Printf("AuthorizeEndpoint: User roles and APIs:")
			for _, role := range user.Roles {
				log.Printf("AuthorizeEndpoint: Role: %s (ID: %d)", role.Name, role.ID)
				for _, api := range role.Apis {
					log.Printf("AuthorizeEndpoint:   - %s %s", api.Tipo, api.Endpoint)
				}
			}
			http.Error(w, "no tiene permisos para este endpoint", http.StatusForbidden)
			return
		}

		log.Printf("AuthorizeEndpoint: Request authorized successfully")
		next.ServeHTTP(w, r)
	})
}

// Define un tipo propio para la clave de contexto
type contextKey string

const UserContextKey contextKey = "user"
