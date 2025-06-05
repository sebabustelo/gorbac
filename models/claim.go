package models

import jwt "github.com/dgrijalva/jwt-go"

//Claim nos permite crear un token
type Claim struct {
	User `json:"user"`
	jwt.StandardClaims
}
