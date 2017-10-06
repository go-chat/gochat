package store

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/go-chat/gochat/apperror"
	"github.com/go-chat/gochat/model"
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
