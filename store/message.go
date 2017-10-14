package store

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/go-chat/gochat/apperror"
	"github.com/go-chat/gochat/model"
)

type MessageStore struct {
	*Store
}

func NewMessageStore(store *Store) *MessageStore {
	return &MessageStore{store}
}

func (gs *MessageStore) Save(message *model.Message) *apperror.AppError {
	if err := gs.SQLStore.Create(message).Error; err != nil {
		logrus.Errorf("cannot store message, err = %v", err)
		return apperror.NewAppError(err, "cannot store message", http.StatusInternalServerError)
	}

	return nil
}
