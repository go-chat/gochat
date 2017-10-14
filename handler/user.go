package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/go-chat/gochat/helper"
	"github.com/go-chat/gochat/model"
)

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
	user.Salt = helper.NewSalt()
	user.Password = helper.HashPassword(p.Password, user.Salt)

	apperr := Srv.Store.IUser.Save(user)
	if apperr != nil {
		logrus.WithFields(lf).Errorf("failed to save user, err = %v", apperr)
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

func Login(w http.ResponseWriter, r *http.Request) {
	lf := logrus.Fields{"func": "handler.Login"}

	type LoginRequest struct {
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"required"`
	}

	var p = &LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		logrus.WithFields(lf).Errorf("Failed to decode body %v", err)
		encodeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	valid, err := govalidator.ValidateStruct(p)
	if err != nil {
		logrus.WithFields(lf).Errorf("Failed to validate body, err = %v", err)
		encodeErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if !valid {
		logrus.WithFields(lf).Errorf("invalid params, err = %v", err)
		encodeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, apperr := Srv.Store.IUser.Login(p.Email, p.Password)
	if apperr != nil {
		logrus.WithFields(lf).Error("failed to login")
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
		Token string `json:"token"`
	}{Token: token})
}
