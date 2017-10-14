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
	IGroup   IGroup
	IMessage IMessage
}

type IUser interface {
	Save(user *model.User) *apperror.AppError
	Login(email, password string) (*model.User, *apperror.AppError)
	GetUserFromToken(token string) (*model.User, *apperror.AppError)
	GetUserFromEmail(email string) (*model.User, *apperror.AppError)
}

type IGroup interface {
	GetGroupsByUserID(userID int64) ([]model.Group, *apperror.AppError)
	Save(group *model.Group) *apperror.AppError
	AddMembersToGroupChat(userIDs []int64, groupID int64) *apperror.AppError
	RemoveMembersToGroupChat(userIDs []int64, groupID int64) *apperror.AppError
	DeleteGroup(userID, groupID int64) *apperror.AppError
	ListMessagesOfGroup(groupID int64, limit, offset int) ([]model.Message, *apperror.AppError)
}

type IMessage interface {
	Save(msg *model.Message) *apperror.AppError
}

func NewStore(config *config.Config) *Store {
	store := &Store{}
	store.SQLStore = getSQLStore(config.SQL.DriverName, config.SQL.DataSource)
	store.SQLStore.LogMode(true)
	migrate(store.SQLStore)

	store.IUser = NewUserStore(store)
	store.IGroup = NewGroupStore(store)
	store.IMessage = NewMessageStore(store)

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
	db.AutoMigrate(&model.Message{})
	db.AutoMigrate(&model.Group{})
	db.AutoMigrate(&model.GroupUser{})
}
