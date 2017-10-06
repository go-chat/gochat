package handler

import (
	"errors"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chat/gochat/model"
)

var (
	TokenInvalid = "token is invalid"
)

// CORS serves the cross origin
func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// TokenAuthMiddleware is auth middlware. Set this header in your request to get here.
// Authorization: Bearer `token`
func TokenAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("RUN_MODE") == "dev" {
			next.ServeHTTP(w, r)
			return
		}

		bearerToken := r.Header.Get("Authorization")
		if len(bearerToken) < 7 {
			encodeErrorResponse(w, errors.New(TokenInvalid), http.StatusForbidden)
			return
		}

		if bearerToken != "" && len(bearerToken) >= 7 {
			if bearerToken[0:7] != "Bearer " {
				encodeErrorResponse(w, errors.New(TokenInvalid), http.StatusForbidden)
				return
			}
		}

		token := bearerToken[7:]

		parseToken, err := jwt.ParseWithClaims(token, &model.CustomClaims{}, func(_token *jwt.Token) (interface{}, error) {
			b := ([]byte(os.Getenv("KHUYA_SECRET")))
			return b, nil

		})
		if err != nil {
			encodeErrorResponse(w, errors.New(TokenInvalid), http.StatusForbidden)
			return
		}

		claims, ok := parseToken.Claims.(*model.CustomClaims)
		if !ok {
			encodeErrorResponse(w, errors.New(TokenInvalid), http.StatusForbidden)
			return
		}

		err = claims.Valid()
		if err != nil {
			encodeErrorResponse(w, errors.New(TokenInvalid), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
