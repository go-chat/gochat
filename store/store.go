package store

import (
	"time"

	"github.com/apex/log"
	"github.com/go-chat/gochat/apperror"
	"github.com/go-chat/gochat/config"
	"github.com/go-chat/gochat/model"
	"github.com/jinzhu/gorm"
)

type Store struct {
	SQLStore *gorm.DB
	IUser    IUser
}

type IUser interface {
	Save(user *model.User) *apperror.AppError
}

func NewStore(config *config.Config) *Store {
	store := &Store{}
	store.SQLStore = getSQLStore(config.SQL.DriverName, config.SQL.DataSource)
	store.SQLStore.LogMode(true)
	migrate(store.SQLStore)

	store.IUser = NewUserStore(store)

	return store
}

func getSQLStore(driver, dataSource string) *gorm.DB {
	db, err := gorm.Open(driver, dataSource)
	if err != nil {
		log.Errorf("Failed to open sql connection to err: %v", err)
		time.Sleep(time.Second)
		panic("Failed to open sql connection" + err.Error())
	}

	return db
}

func (store *Store) Close() {
	store.SQLStore.Close()
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
}
