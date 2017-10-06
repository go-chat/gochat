package model

import jwt "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	Avatar             string `json:"avatar"`
	jwt.StandardClaims `json:"-"`
}
