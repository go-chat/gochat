package handler

import (
	"net/http"

	"github.com/go-chat/gochat/model"
)

func Index(w http.ResponseWriter, r *http.Request) {
	encodeSuccessResponse(w, "hehe")
}

func Register(w http.ResponseWriter, r *http.Request) {
	user := &model.User{
		Name:     "hehe",
		Email:    "lnthach2110@gmail.com",
		Password: "123",
	}
	apperr := Srv.Store.IUser.Save(user)
	if apperr != nil {
		encodeAppErrorResponse(w, apperr)
		return
	}

	encodeSuccessResponse(w, "user created successfully")
}
