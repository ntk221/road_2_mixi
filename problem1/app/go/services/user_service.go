package service

import (
	"database/sql"
	"problem1/model"
	"problem1/repository"
	"problem1/types"
)

type UserService interface {
	GetFriendList(user_id int) ([]model.User, error)
	GetFriendListFromUsers([]model.User) ([]model.User, error)
	GetFriendListWithPagenation(user_id int, params types.PagenationParams) ([]model.User, error)
	GetFriendListFromUsersWithPagenation([]model.User, types.PagenationParams) ([]model.User, error)
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

func (us UserServiceImpl) GetFriendListWithPagenation(user_id int, params types.PagenationParams) ([]model.User, error) {
	friends, err := us.ur.GetFriendsByID(user_id, params, us.db)
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

func (us UserServiceImpl) GetFriendListFromUsersWithPagenation(friendList []model.User, params types.PagenationParams) ([]model.User, error) {
	fofs := make([]model.User, 0)

	for _, friend := range friendList {
		fof, err := us.GetFriendList(friend.UserID)
		if err != nil {
			return nil, err
		}
		fofs = append(fofs, fof...)
	}

	fofs = pagenate(params, fofs)
	return fofs, nil
}

func (us UserServiceImpl) GetFriendList(user_id int) ([]model.User, error) {
	params := types.PagenationParams{
		Offset: 0,
		Limit:  10,
	}

	return us.GetFriendListWithPagenation(user_id, params)
}

func (us UserServiceImpl) GetFriendListFromUsers(friendList []model.User) ([]model.User, error) {
	params := types.PagenationParams{
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

func pagenate(params types.PagenationParams, users []model.User) []model.User {
	if params.Offset > len(users) {
		return make([]model.User, 0)
	}

	if params.Offset+params.Limit > len(users) {
		return users[params.Offset:]
	}

	return users[params.Offset : params.Offset+params.Limit]
}
