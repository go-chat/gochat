package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chat/gochat/apperror"
	"github.com/go-chat/gochat/model"
)

func generateToken(user *model.User) (string, *apperror.AppError) {
	lf := logrus.Fields{"func": "handler.generateToken"}

	tokenExpireDays, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRE_DAYS"))
	if err != nil {
		logrus.WithFields(lf).Errorf("failed to parse token expire days, err =%v", err)
		return "", apperror.NewAppError(err, "cannot parse token expire day", http.StatusInternalServerError)
	}

	claims := &model.CustomClaims{
		user.ID,
		user.Email,
		user.Name,
		user.Avatar,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24*tokenExpireDays)).Unix(),
			Issuer:    "khuya",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(os.Getenv("KHUYA_SECRET")))
	if err != nil {
		logrus.WithFields(lf).Errorf("failed to get signed string, err=%v", err)
		return "", apperror.NewAppError(err, "cannot sign string", http.StatusInternalServerError)
	}

	return tokenString, nil
}

func getUserFromToken(r *http.Request) (*model.User, *apperror.AppError) {
	lf := logrus.Fields{"func": "handler.getUserFromToken"}

	token := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(token) < 7 {
		logrus.WithFields(lf).Error("token is invalid")
		return nil, apperror.NewAppError(errors.New("token is invalid"), "token is invalid", http.StatusBadRequest)
	}
	token = token[7:]

	return Srv.Store.IUser.GetUserFromToken(token)
}
