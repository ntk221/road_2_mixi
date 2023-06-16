package service

import (
	"database/sql"
	"problem1/model"
	"problem1/repository"
)

type UserService interface {
	GetFriendList(user_id int) ([]model.User, error)
	GetFriendListFromUsers([]model.User) ([]model.User, error)
}

type UserServiceImpl struct {
	db *sql.DB
	ur *repository.UserRepositoryImpl
}

func NewUserService(db *sql.DB, ur *repository.UserRepositoryImpl) UserService {
	return &UserServiceImpl{
		db: db,
		ur: ur,
	}
}

func (us UserServiceImpl) GetFriendList(user_id int) ([]model.User, error) {
	friends, err := us.ur.GetFriendsByID(user_id, us.db)
	if err != nil {
		return nil, err
	}

	blockedUsers, err := us.ur.GetBlockedUsersByID(user_id, us.db)
	if err != nil {
		return nil, err
	}

	filteredFriends := fileterBlockedFriends(friends, blockedUsers)

	return filteredFriends, nil
}

func (us UserServiceImpl) GetFriendListFromUsers(friendList []model.User) ([]model.User, error) {
	fofs := make([]model.User, 0)

	for _, friend := range friendList {
		fof, err := us.GetFriendList(friend.UserID)
		if err != nil {
			return nil, err
		}
		fofs = append(fofs, fof...)
	}
	return fofs, nil
}

func fileterBlockedFriends(friends []model.User, blocked []model.User) []model.User {
	friendNames := make([]model.User, 0)

	for _, friend := range friends {
		if !contains(blocked, friend) {
			friendNames = append(friendNames, friend)
		}
	}

	return friendNames
}

func contains(slice []model.User, value model.User) bool {
	for _, item := range slice {
		if item.ID == value.ID {
			return true
		}
	}
	return false
}
