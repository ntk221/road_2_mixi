package repository

import (
	"database/sql"
	"fmt"
	"problem1/model"
	"strings"
)

type UserRepository interface {
	GetFriendsByID(user_id int, db Queryer) ([]model.User, error)
	GetBlockedUsersByID(user_id int, db Queryer) ([]model.User, error)
	GetByID(user_id int, db Queryer) (model.User, error)
	// GetFriendNames(ids []string) ([]string, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (ur *UserRepositoryImpl) GetFriendsByID(user_id int, db Queryer) ([]model.User, error) {
	query := `SELECT user1_id, user2_id FROM friend_link WHERE user1_id = ? OR user2_id = ?`
	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendIDs []int
	for rows.Next() {
		var friend1ID, friend2ID int
		if err := rows.Scan(&friend1ID, &friend2ID); err != nil {
			return nil, err
		}
		if friend1ID == user_id {
			friendIDs = append(friendIDs, friend2ID)
		} else {
			friendIDs = append(friendIDs, friend1ID)
		}
	}

	friends := make([]model.User, 0)
	for _, friendID := range friendIDs {
		friend, err := ur.GetByID(friendID, db)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (ur *UserRepositoryImpl) GetBlockedUsersByID(user_id int, db Queryer) ([]model.User, error) {
	query := `SELECT user1_id, user2_id FROM block_list WHERE user1_id = ? OR user2_id = ?`
	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockedIDs []int
	for rows.Next() {
		var user1ID, user2ID int
		if err := rows.Scan(&user1ID, &user2ID); err != nil {
			return nil, err
		}
		if user1ID != user_id {
			blockedIDs = append(blockedIDs, user1ID)
		}
		if user2ID != user_id {
			blockedIDs = append(blockedIDs, user2ID)
		}
	}

	blocked := make([]model.User, 0)
	for _, blockedID := range blockedIDs {
		blockedUser, err := ur.GetByID(blockedID, db)
		if err != nil {
			return nil, err
		}
		blocked = append(blocked, blockedUser)
	}

	return blocked, nil
}

func (ur *UserRepositoryImpl) GetByID(user_id int, db Queryer) (model.User, error) {
	query := `SELECT id, user_id, name FROM users WHERE user_id = ?`
	row := db.QueryRow(query, user_id)

	var user model.User
	if err := row.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, fmt.Errorf("user not found")
		}
		return model.User{}, err
	}
	return user, nil
}

func replacePlaceholders(query string, argCount int) string {
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = "?"
	}
	return strings.Replace(query, "?", strings.Join(placeholders, ","), -1)
}
