package usecases

import (
	"problem1/domain"
	"problem1/infra"
)

type UserService interface {
	GetFriendList(user_id domain.UserID) ([]domain.User, error)
	GetFriendListFromUsers([]domain.User, int) ([]domain.User, error)
	GetUserByID(user_id domain.UserID) (domain.User, error)
}

type UserServiceImpl struct {
	db domain.Database
	// ur UserRepository
}

func NewUserService(db domain.Database) UserService {
	return &UserServiceImpl{
		db: db,
		// ur: ur,
	}
}

func (us UserServiceImpl) GetUserByID(user_id domain.UserID) (domain.User, error) {
	ur := infra.NewUserRepository()
	tx, err := us.db.Begin()
	if err != nil {
		return domain.User{}, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		}
	}()

	user, err := ur.GetByID(user_id, tx)
	if err != nil {
		return domain.User{}, err
	}

	err = tx.Commit()
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (us UserServiceImpl) GetFriendList(user_id domain.UserID) ([]domain.User, error) {
	ur := infra.NewUserRepository()
	tx, err := us.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		}
	}()

	user, err := ur.GetByID(user_id, tx)
	if err != nil {
		return nil, err
	}

	friendIDs := user.GetFriendList()

	friends, err := us.getUsersByIDs(friendIDs)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// depth の分だけ再帰呼び出しを行う
// そのたびにforループを呼び出している
// これは，depthが大きくなると，再帰呼び出しの回数が増えるため，効率が悪い
// TODO: 計算量を改善する
func (us UserServiceImpl) GetFriendListFromUsers(userList []domain.User, depth int) ([]domain.User, error) {
	friends := make([]domain.User, 0)

	for _, user := range userList {
		v, err := us.GetFriendList(user.UserID)
		if err != nil {
			return nil, err
		}
		friends = append(friends, v...)

		if depth > 1 {
			// 友人の友人のリストを再帰的に取得
			recursiveFriends, err := us.GetFriendListFromUsers(v, depth-1)
			if err != nil {
				return nil, err
			}
			friends = append(friends, recursiveFriends...)
		}
	}

	return friends, nil
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
