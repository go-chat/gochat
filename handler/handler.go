package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chat/gochat/apperror"
	"github.com/go-chat/gochat/model"
)

func Index(w http.ResponseWriter, r *http.Request) {
	encodeSuccessResponse(w, "hehe")
}

func Register(w http.ResponseWriter, r *http.Request) {
	lf := logrus.Fields{"func": "handler.Register"}

	type RegisterRequest struct {
		Email           string `json:"email" valid:"email"`
		Name            string `json:"name" valid:"required"`
		Password        string `json:"password" valid:"required"`
		ConfirmPassword string `json:"confirm_password" valid:"required"`
	}

	var p = &RegisterRequest{}

	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		logrus.WithFields(lf).Errorf("failed to decode body, err = %v", err)
		encodeErrorResponse(w, fmt.Errorf("failed to decode body, err = %v", err), http.StatusBadRequest)
		return
	}

	valid, err := govalidator.ValidateStruct(p)
	if err != nil {
		logrus.WithFields(lf).Errorf("failed to validate body, err =  %v", err)
		encodeErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if !valid {
		logrus.WithFields(lf).Errorf("invalid params, err =  %v", err)
		encodeErrorResponse(w, errors.New("struct is not valid"), http.StatusBadRequest)
		return
	}

	if p.Password != p.ConfirmPassword {
		logrus.WithFields(lf).Errorf("confirm password does not match, err =  %v", err)
		encodeErrorResponse(w, errors.New("confirm password does not match"), http.StatusBadRequest)
		return
	}

	user := &model.User{}
	user.Email = p.Email
	user.Name = p.Name
	user.Salt = NewSalt()
	user.Password = HashPassword(user.Password, user.Salt)

	apperr := Srv.Store.IUser.Save(user)
	if apperr != nil {
		encodeAppErrorResponse(w, apperr)
		return
	}

	token, appErr := generateToken(user)
	if appErr != nil {
		logrus.WithFields(lf).Error("failed to generate token")
		encodeAppErrorResponse(w, appErr)
		return
	}

	encodeSuccessResponse(w, struct {
		Token  string `json:"token"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}{Token: token,
		Name:   user.Name,
		Email:  user.Email,
		Avatar: user.Avatar})
}

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
