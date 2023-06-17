package service

import (
	"database/sql"
	"problem1/model"
	"problem1/repository"
)

type UserService interface {
	GetFriendList(user_id int) ([]model.User, error)
	GetFriendListFromUsers([]model.User) ([]model.User, error)
	GetFriendListWithPagenation(user_id int, params PagenationParams) ([]model.User, error)
	GetFriendListFromUsersWithPagenation([]model.User, PagenationParams) ([]model.User, error)
}

type UserServiceImpl struct {
	db *sql.DB
	ur *repository.UserRepositoryImpl
}

type PagenationParams struct {
	Offset int
	Limit  int
}

func NewUserService(db *sql.DB, ur *repository.UserRepositoryImpl) UserService {
	return &UserServiceImpl{
		db: db,
		ur: ur,
	}
}

func (us UserServiceImpl) GetFriendListWithPagenation(user_id int, params PagenationParams) ([]model.User, error) {
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

func (us UserServiceImpl) GetFriendListFromUsersWithPagenation(friendList []model.User, params PagenationParams) ([]model.User, error) {
	fofs := make([]model.User, 0)

	for _, friend := range friendList {
		fof, err := us.GetFriendListWithPagenation(friend.UserID, params)
		if err != nil {
			return nil, err
		}
		fofs = append(fofs, fof...)
	}
	return fofs, nil
}

func (us UserServiceImpl) GetFriendList(user_id int) ([]model.User, error) {
	params := PagenationParams{
		Offset: 0,
		Limit:  10,
	}

	return us.GetFriendListWithPagenation(user_id, params)
}

func (us UserServiceImpl) GetFriendListFromUsers(friendList []model.User) ([]model.User, error) {
	params := PagenationParams{
		Offset: 0,
		Limit:  10,
	}

	return us.GetFriendListFromUsersWithPagenation(friendList, params)
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
