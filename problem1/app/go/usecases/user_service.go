package usecases

import (
	"database/sql"
	"problem1/domain"
	"problem1/infra"
)

type UserService interface {
	GetFriendList(user_id int) ([]domain.User, error)
	GetFriendListFromUsers([]domain.User) ([]domain.User, error)
	// GetFriendListWithPagenation(user_id int, params types.PagenationParams) ([]model.User, error)
	// GetFriendListFromUsersWithPagenation([]model.User, types.PagenationParams) ([]model.User, error)
}

type UserServiceImpl struct {
	db infra.Database
	ur UserGetter
}

func NewUserService(db *sql.DB, ur *infra.UserRepositoryImpl) UserService {
	return &UserServiceImpl{
		db: db,
		ur: ur,
	}
}

func (us UserServiceImpl) GetFriendList(user_id int) ([]domain.User, error) {
	user, err := us.ur.GetByID(user_id, us.db)
	if err != nil {
		return nil, err
	}

	realFriends, err := us.getRealFriends(user)
	if err != nil {
		return nil, err
	}

	filteredFriends, err := us.filterWithBlockLink(user, realFriends)
	if err != nil {
		return nil, err
	}

	return filteredFriends, nil
}

func (us UserServiceImpl) filterWithBlockLink(user domain.User, friends []domain.User) ([]domain.User, error) {
	filteredFriends := make([]domain.User, 0)

	for _, friend := range friends {
		if !contains(user.BlockList, friend.UserID) {
			filteredFriends = append(filteredFriends, friend)
		}
	}

	return filteredFriends, nil
}

func (us UserServiceImpl) getRealFriends(user domain.User) ([]domain.User, error) {
	realFriends := make([]domain.User, 0)

	friendIDs, err := user.GetFriendList()
	if err != nil {
		return nil, err
	}

	for _, friendID := range friendIDs {
		friend, err := us.ur.GetByID(friendID, us.db)
		if err != nil {
			return nil, err
		}
		if contains(friend.FriendList, user.UserID) && !contains(user.BlockList, friend.UserID) {
			realFriends = append(realFriends, friend)
		}
	}

	return realFriends, nil
}

func (us UserServiceImpl) GetFriendListFromUsers(friendList []domain.User) ([]domain.User, error) {
	fofs := make([]domain.User, 0)

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
