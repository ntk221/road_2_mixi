package service

import (
	"database/sql"
	"fmt"
	"problem1/model"
	"problem1/repository"
)

type UserService interface {
	GetFriendList(user_id int) ([]model.User, error)
	GetFriendListFromUsers([]model.User) ([]model.User, error)
	// GetFriendListWithPagenation(user_id int, params types.PagenationParams) ([]model.User, error)
	// GetFriendListFromUsersWithPagenation([]model.User, types.PagenationParams) ([]model.User, error)
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
	user, err := us.ur.GetByID(user_id, us.db)
	if err != nil {
		return nil, err
	}

	// user が 友人だと思っているだけではなく，相手もuserを友人だと思っていることが条件
	realFriends, err := us.getRealFriends(user)
	if err != nil {
		return nil, err
	}

	fmt.Println(realFriends)

	filteredFriends, err := us.fileterWithBlockLink(user, realFriends)
	if err != nil {
		return nil, err
	}

	return filteredFriends, nil
}

func (us UserServiceImpl) fileterWithBlockLink(user model.User, friends []model.User) ([]model.User, error) {
	filteredFriends := make([]model.User, 0)

	for _, friend := range friends {
		if !contains(user.BlockList, friend.UserID) {
			filteredFriends = append(filteredFriends, friend)
		}
	}

	return filteredFriends, nil
}

func (us UserServiceImpl) getRealFriends(user model.User) ([]model.User, error) {
	realFriends := make([]model.User, 0)

	for _, friendID := range user.FriendList {
		friend, err := us.ur.GetByID(friendID, us.db)
		if err != nil {
			return nil, err
		}

		fmt.Println(friend)

		if contains(friend.FriendList, user.UserID) && !contains(user.BlockList, friend.UserID) {
			realFriends = append(realFriends, friend)
		}
	}

	return realFriends, nil
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

	// fofs = pagenate(params, fofs)
	return fofs, nil
}

func contains(slice []int, value int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
