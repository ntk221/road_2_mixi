package usecases

import (
	"problem1/domain"
)

type UserService interface {
	GetFriendList(user_id domain.UserID) ([]domain.User, error)
	GetFriendListFromUsers([]domain.User) ([]domain.User, error)
	GetUserByID(user_id domain.UserID) (domain.User, error)
	// GetUsersByIDs() ([]domain.User, error)
	// GetFriendListWithPagenation(user_id int, params types.PagenationParams) ([]model.User, error)
	// GetFriendListFromUsersWithPagenation([]model.User, types.PagenationParams) ([]model.User, error)
}

type UserServiceImpl struct {
	qx domain.QueryerTx
	ur UserRepository
}

func NewUserService(qx domain.QueryerTx, ur UserGetter) UserService {
	return &UserServiceImpl{
		qx: qx,
		ur: ur,
	}
}

func (us UserServiceImpl) GetUserByID(user_id domain.UserID) (domain.User, error) {
	user, err := us.ur.GetByID(user_id, us.qx)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (us UserServiceImpl) GetFriendList(user_id domain.UserID) ([]domain.User, error) {
	user, err := us.ur.GetByID(user_id, us.qx)
	if err != nil {
		return nil, err
	}

	friendIDs := user.GetFriendList()

	friends, err := us.getUsersByIDs(friendIDs)
	if err != nil {
		return nil, err
	}

	return friends, nil
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

func (us UserServiceImpl) getUsersByIDs(user_ids []domain.UserID) ([]domain.User, error) {
	users := make([]domain.User, 0)
	for _, user_id := range user_ids {
		user, err := us.GetUserByID(user_id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
