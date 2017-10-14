package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/go-chat/gochat/model"
	"github.com/gorilla/mux"
)

func GetGroupsByUserID(w http.ResponseWriter, r *http.Request) {
	lf := logrus.Fields{"func": "handler.GetGroupsByUserID"}

	user, appErr := getUserFromToken(r)
	if appErr != nil {
		logrus.WithFields(lf).WithError(appErr.Error).Error("cannot get user from token")
		encodeAppErrorResponse(w, appErr)
		return
	}

	groups, appErr := Srv.Store.IGroup.GetGroupsByUserID(user.ID)
	if appErr != nil {
		logrus.WithFields(lf).WithError(appErr.Error).Error("cannot get groups by user id")
		encodeAppErrorResponse(w, appErr)
		return
	}

	encodeSuccessResponse(w, groups)
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	lf := logrus.Fields{"func": "handler.CreateGroup"}

	type CreateGroupRequest struct {
		Name string `json:"name" valid:"required"`
		Meta string `json:"meta"`
	}

	var p = &CreateGroupRequest{}

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

	group := &model.Group{}
	group.Name = p.Name
	group.Meta = p.Meta

	apperr := Srv.Store.IGroup.Save(group)
	if apperr != nil {
		logrus.WithFields(lf).Errorf("failed to save group, err = %v", apperr)
		encodeAppErrorResponse(w, apperr)
		return
	}

	encodeSuccessResponse(w, map[string]interface{}{"success": true})
}

func AddMembersToGroupChat(w http.ResponseWriter, r *http.Request) {
	lf := logrus.Fields{"func": "handler.AddMembersToGroupChat"}

	type AddMembersToGroupChatRequest struct {
		UserIDs []int64 `json:"user_ids" validate:"required"`
		GroupID int64   `json:"group_id" validate:"required"`
	}

	var p = &AddMembersToGroupChatRequest{}

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

	apperr := Srv.Store.IGroup.AddMembersToGroupChat(p.UserIDs, p.GroupID)
	if apperr != nil {
		logrus.WithFields(lf).Errorf("failed to add member to group, err = %v", apperr)
		encodeAppErrorResponse(w, apperr)
		return
	}

	encodeSuccessResponse(w, map[string]interface{}{"success": true})
}

func RemoveMembersToGroupChat(w http.ResponseWriter, r *http.Request) {
	lf := logrus.Fields{"func": "handler.RemoveMembersToGroupChat"}

	type RemoveMembersToGroupChatRequest struct {
		UserIDs []int64 `json:"user_ids" validate:"required"`
		GroupID int64   `json:"group_id" validate:"required"`
	}

	var p = &RemoveMembersToGroupChatRequest{}

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

	appErr := Srv.Store.IGroup.RemoveMembersToGroupChat(p.UserIDs, p.GroupID)
	if appErr != nil {
		logrus.WithFields(lf).Errorf("failed to add member to group, err = %v", appErr)
		encodeAppErrorResponse(w, appErr)
		return
	}

	encodeSuccessResponse(w, map[string]interface{}{"success": true})
}

func DeleteGroupChat(w http.ResponseWriter, r *http.Request) {
	lf := logrus.Fields{"func": "handler.DeleteGroupChat"}

	vars := mux.Vars(r)
	idstr := vars["id"]

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		logrus.WithFields(lf).Errorf("failed to parse int from param, err = %v", err)
		encodeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, appErr := getUserFromToken(r)
	if appErr != nil {
		logrus.WithFields(lf).WithError(appErr.Error).Error("cannot get user from token")
		encodeAppErrorResponse(w, appErr)
		return
	}

	appErr = Srv.Store.IGroup.DeleteGroup(user.ID, id)
	if appErr != nil {
		logrus.WithFields(lf).Errorf("failed to add member to group, err = %v", appErr)
		encodeAppErrorResponse(w, appErr)
		return
	}

	encodeSuccessResponse(w, map[string]interface{}{"success": true})
}

func ListMessagesOfGroup(w http.ResponseWriter, r *http.Request) {
	lf := logrus.Fields{"func": "handler.ListMessagesOfGroup"}

	slimit := r.URL.Query().Get("limit")
	soffset := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(slimit)
	if err != nil {
		logrus.WithFields(lf).Errorf("failed to parse int limit, limit = %s, err = %v", slimit, err)
		encodeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(soffset)
	if err != nil {
		logrus.WithFields(lf).Errorf("failed to parse int offset, offset = %s, err = %v", soffset, err)
		encodeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	idstr := vars["id"]

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		logrus.WithFields(lf).Errorf("failed to parse int from param, err = %v", err)
		encodeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	messages, appErr := Srv.Store.IGroup.ListMessagesOfGroup(id, limit, offset)
	if appErr != nil {
		logrus.WithFields(lf).WithError(appErr.Error).Error("cannot get messages by group id")
		encodeAppErrorResponse(w, appErr)
		return
	}

	encodeSuccessResponse(w, messages)
}
