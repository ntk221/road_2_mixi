package usecases

import (
	"problem1/domain/entity"
	"problem1/domain/repository"
	"problem1/infra"

	"log"
)

type UserService interface {
	GetFriendList(user_id entity.UserID) (*entity.UserCollection, error)
	GetFriendListFromUsers(*entity.UserCollection, int) (*entity.UserCollection, error)
	GetUserByID(user_id entity.UserID) (*entity.User, error)
}

type UserServiceImpl struct {
	db repository.Database
	// ur UserRepository
}

func NewUserService(db repository.Database) UserService {
	return &UserServiceImpl{
		db: db,
		// ur: ur,
	}
}

func (us UserServiceImpl) GetUserByID(user_id entity.UserID) (*entity.User, error) {
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

	log.Println("user: ", user)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// id に対応するユーザーの友達のユーザー情報を取得する
// ユーザー情報はuniqueにする
func (us UserServiceImpl) GetFriendList(user_id entity.UserID) (*entity.UserCollection, error) {
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
func (us UserServiceImpl) GetFriendListFromUsers(userList *entity.UserCollection, depth int) (*entity.UserCollection, error) {
	friends := entity.NewUserCollection([]*entity.User{})

	for _, user := range userList.Users {
		v, err := us.GetFriendList(user.UserID)
		if err != nil {
			return nil, err
		}
		friends = friends.Merge(v)

		if depth > 1 {
			// 友人の友人のリストを再帰的に取得
			recursiveFriends, err := us.GetFriendListFromUsers(v, depth-1)
			if err != nil {
				return nil, err
			}
			friends = friends.Merge(recursiveFriends)
		}
	}

	return friends, nil
}

func (us UserServiceImpl) getUsersByIDs(user_ids []entity.UserID) (*entity.UserCollection, error) {
	users := make([]*entity.User, 0)
	for _, user_id := range user_ids {
		user, err := us.GetUserByID(user_id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	userCollection := entity.NewUserCollection(users)
	userCollection = userCollection.GetUniqueUsers()
	return userCollection, nil
}
