package authentication

import (
	"api-rbac/db"
	"api-rbac/models"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	responses "api-rbac/helpers"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	//privateBytes string
)

func init() {

	privateBytes, err := ioutil.ReadFile("./private.rsa")
	if err != nil {
		log.Fatal("No se puede leer el arhivo privado")

	}

	publicBytes, err := ioutil.ReadFile("./public.rsa.pub")
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
		responses.ERROR(w, http.StatusBadRequest, errors.New("datos de login inválidos"))
		return
	}

	// 1. Verifica el token de Google
	payload, err := idtoken.Validate(r.Context(), req.IdToken, "285580531726-v8385mgjrgebb9arkm0hr9beqolvrrm4.apps.googleusercontent.com")
	if err != nil {
		fmt.Println("Error validando token de Google:", err)
		responses.ERROR(w, http.StatusUnauthorized, errors.New("token de Google inválido"))
		return
	}

	// 2. Extrae datos del usuario de Google
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	sub, _ := payload.Claims["sub"].(string) // Google user ID

	// 3. Busca o crea el usuario en tu base de datos
	var user models.User
	db := db.Instance()
	err = db.Where("email = ?", email).First(&user).Preload("Roles").Error
	if err != nil {
		// Si no existe, créalo
		user = models.User{
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
		fmt.Println("ususario a crear ", user)
		db.Create(&user)
	}

	//si no existe el usuario, se crea
	if user.ID == 0 {
		user = models.User{
			Email:      email,
			Name:       name,
			Provider:   "google",
			ProviderID: sub,
			Active:     true,
		}
		if err := db.Create(&user).Error; err != nil {
			responses.ERROR(w, http.StatusInternalServerError, errors.New("error al crear el usuario"))
			return
		}
	}

	// 4. Genera tu JWT y responde
	user.Password = ""
	user.Token = GenerateJWT(user)
	responses.JSON(w, http.StatusOK, user)
}

// RefreshToken Genera un nuevo token a partir de un token valido y lo entrega
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenValidation := TokenValid(r)

	if tokenValidation {
		tokenString := strings.Split(r.Header["Authorization"][0], "Bearer ")
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString[1], jwt.MapClaims{})
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
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
		token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Su token no es válido")
			return
		}

		// Extrae el claim y lo guarda en el contexto
		if claims, ok := token.Claims.(*models.Claim); ok {
			ctx := context.WithValue(r.Context(), userContextKey, *claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

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
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				log.Println(err)
				//fmt.Fprintln(w, "Su token ha expirado")
				return false
			case jwt.ValidationErrorClaimsInvalid:
				//fmt.Fprintln(w, "La firma del token no coincide")
				return false
			default:
				//fmt.Fprintln(w, "Su token no es valido 1")
				return false
			}
		default:
			//fmt.Fprintln(w, "Su token no es válido 2")
			return
		}
	}

	if token.Valid {
		//w.WriteHeader(http.StatusAccepted)
		//fmt.Fprintln(w, "Bienvenido al sistema")
		status = true
	} else {
		//w.WriteHeader(http.StatusUnauthorized)
		//fmt.Fprintln(w, "Su token no es válido")
		status = false
	}

	return status

}

// AuthorizeEndpoint verifica si el usuario tiene permiso para acceder al endpoint
func AuthorizeEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim, ok := r.Context().Value(userContextKey).(models.Claim)
		if !ok {
			http.Error(w, "no autorizado", http.StatusUnauthorized)
			return
		}
		userID := claim.User.ID
		endpoint := r.URL.Path

		db := db.Instance()

		var user models.User
		err := db.Where("id = ?", userID).
			Preload("Roles").
			Preload("Roles.Apis").
			First(&user).Error
		if err != nil {
			http.Error(w, "usuario no encontrado", http.StatusUnauthorized)
			return
		}

		hasPermission := false
		for _, role := range user.Roles {
			for _, api := range role.Apis {
				if api.Endpoint == endpoint {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			http.Error(w, "no tiene permisos para este endpoint", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Define un tipo propio para la clave de contexto
type contextKey string

const userContextKey contextKey = "user"
