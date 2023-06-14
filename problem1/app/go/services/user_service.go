package service

import (
	"problem1/model"
	"problem1/repository"
)

type UserService interface {
	GetFriendList(id string) ([]model.User, error)
}

type UserServiceImpl struct {
	ur repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &UserServiceImpl{ur: ur}
}

func (us UserServiceImpl) GetFriendList(id string) ([]model.User, error) {
	friends, err := us.ur.GetFriendsByID(id)
	if err != nil {
		return nil, err
	}

	blockedUsers, err := us.ur.GetBlockedUsersByID(id)
	if err != nil {
		return nil, err
	}

	filteredFriends := fileterBlockedFriends(friends, blockedUsers)

	return filteredFriends, nil
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