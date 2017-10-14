package store

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chat/gochat/apperror"
	"github.com/go-chat/gochat/helper"
	"github.com/go-chat/gochat/model"
	"github.com/jinzhu/gorm"
)

type UserStore struct {
	*Store
}

func NewUserStore(store *Store) *UserStore {
	return &UserStore{store}
}

func (us *UserStore) Save(user *model.User) *apperror.AppError {
	if err := us.SQLStore.Create(user).Error; err != nil {
		logrus.Errorf("cannot store user, err = %v", err)
		return apperror.NewAppError(err, "cannot store user", http.StatusInternalServerError)
	}

	return nil
}

func (us *UserStore) Login(email, password string) (*model.User, *apperror.AppError) {
	user, err := us.GetUserFromEmail(email)
	if err != nil {
		return nil, err
	}

	if !checkUserPass(password, user.Salt, user.Password) {
		return nil, apperror.NewAppError(fmt.Errorf("error or password does not match"), "email or password does not match", http.StatusBadRequest)
	}

	return user, nil
}

func (us *UserStore) GetUserFromEmail(email string) (*model.User, *apperror.AppError) {
	user := &model.User{}
	if err := us.SQLStore.Where("email = ?", email).First(&user).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, apperror.NewAppError(err, "email not found", http.StatusNotFound)
		}

		logrus.Errorf("cannot get user by email, email = %s, err = %v", email, err)
		return nil, apperror.NewAppError(err, "cannot get user by email", http.StatusInternalServerError)
	}

	return user, nil
}

func (us *UserStore) GetUserFromToken(token string) (*model.User, *apperror.AppError) {
	parseToken, err := jwt.ParseWithClaims(token, &model.CustomClaims{}, func(_token *jwt.Token) (interface{}, error) {
		b := ([]byte(os.Getenv("GOCHAT_SECRET")))
		return b, nil
	})
	if err != nil {
		logrus.Errorf("cannot parse token, token = %s, err = %v", token, err)
		return nil, apperror.NewAppError(err, "cannot parse token", http.StatusBadRequest)
	}

	claims := parseToken.Claims.(*model.CustomClaims)
	email := claims.Email

	user, appErr := us.GetUserFromEmail(email)
	if appErr != nil {
		return nil, appErr
	}

	return user, nil
}

func checkUserPass(password, salt, hashedPassword string) bool {
	return helper.HashPassword(password, salt) == hashedPassword
}
