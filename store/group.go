package store

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/go-chat/gochat/apperror"
	"github.com/go-chat/gochat/model"
)

type GroupStore struct {
	*Store
}

func NewGroupStore(store *Store) *GroupStore {
	return &GroupStore{store}
}

func (gs *GroupStore) GetGroupsByUserID(userID int64) ([]model.Group, *apperror.AppError) {
	var groups []model.Group

	if err := gs.SQLStore.Model(&model.Group{}).Joins("join group_users ON group_users.group_id = groups.id").Where("group_users.user_id = ? AND group_users.is_active = ?", userID, true).Find(&groups).Error; err != nil {
		logrus.Errorf("cannot get groups by user id, id = %d, err = %v", userID, err)
		return nil, apperror.NewAppError(err, "cannot get groups by user id", http.StatusInternalServerError)
	}

	return groups, nil
}

func (gs *GroupStore) Save(group *model.Group) *apperror.AppError {
	if err := gs.SQLStore.Create(group).Error; err != nil {
		logrus.Errorf("cannot store group, err = %v", err)
		return apperror.NewAppError(err, "cannot store group", http.StatusInternalServerError)
	}

	return nil
}

func (gs *GroupStore) AddMembersToGroupChat(userIDs []int64, groupID int64) *apperror.AppError {

	for _, v := range userIDs {
		groupUser := &model.GroupUser{
			UserID:   v,
			GroupID:  groupID,
			IsActive: true,
		}

		if err := gs.SQLStore.Create(groupUser).Error; err != nil {
			logrus.Errorf("cannot store group user, err = %v", err)
			return apperror.NewAppError(err, "cannot store group user ", http.StatusInternalServerError)
		}
	}
	return nil
}

func (gs *GroupStore) RemoveMembersToGroupChat(userIDs []int64, groupID int64) *apperror.AppError {

	for _, v := range userIDs {
		if err := gs.SQLStore.Where("user_id = ? AND group_id = ?", v, groupID).Delete(&model.GroupUser{}).Error; err != nil {
			logrus.Errorf("cannot remove group user, err = %v", err)
			return apperror.NewAppError(err, "cannot remove group user ", http.StatusInternalServerError)
		}
	}
	return nil
}

func (gs *GroupStore) DeleteGroup(userID, groupID int64) *apperror.AppError {
	if err := gs.SQLStore.Where("user_id = ? AND id = ?", userID, groupID).Delete(&model.Group{}).Error; err != nil {
		logrus.Errorf("cannot delete group, err = %v", err)
		return apperror.NewAppError(err, "cannot delete group", http.StatusInternalServerError)
	}

	return nil
}

func (gs *GroupStore) ListMessagesOfGroup(groupID int64, limit, offset int) ([]model.Message, *apperror.AppError) {
	var messages []model.Message

	if err := gs.SQLStore.Where("group_id = ?", groupID).Limit(limit).Offset(offset).Find(messages).Error; err != nil {
		logrus.Errorf("cannot get messages of group, id = %d, err = %v", groupID, err)
		return nil, apperror.NewAppError(err, "cannot get messages by group id", http.StatusInternalServerError)
	}

	return messages, nil
}
