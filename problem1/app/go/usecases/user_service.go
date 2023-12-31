package usecases

import (
	"problem1/domain"
)

type UserService interface {
	GetFriendList(user_id int) ([]domain.User, error)
	GetFriendListFromUsers([]domain.User) ([]domain.User, error)
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

func (us UserServiceImpl) GetFriendList(user_id int) ([]domain.User, error) {
	user, err := us.ur.GetByID(user_id, us.qx)
	if err != nil {
		return nil, err
	}

	friends := make([]domain.User, 0)
	friendIDs := user.GetFriendList()
	for _, friendID := range friendIDs {
		friend, err := us.ur.GetByID(friendID, us.qx)
		if err != nil {
			return nil, err
		}
		if friend.IsBlocked(user) || user.IsBlocked(friend) {
			continue
		}
		friends = append(friends, friend)
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
