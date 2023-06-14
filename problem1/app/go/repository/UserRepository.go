package repository

import (
	"database/sql"
	"fmt"
	"problem1/model"
	"strings"
)

type UserRepository interface {
	GetFriendsByID(id string) ([]model.User, error)
	GetBlockedUsersByID(id string) ([]model.User, error)
	GetUser(id string) (string, error)
	// GetFriendNames(ids []string) ([]string, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) GetFriendsByID(id string) ([]model.User, error) {
	query := `SELECT user1_id, user2_id FROM friend_link WHERE user1_id = ? OR user2_id = ?`
	rows, err := ur.db.Query(query, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendIDs []string
	for rows.Next() {
		var friend1ID, friend2ID string
		if err := rows.Scan(&friend1ID, &friend2ID); err != nil {
			return nil, err
		}
		if friend1ID == id {
			friendIDs = append(friendIDs, friend2ID)
		} else {
			friendIDs = append(friendIDs, friend1ID)
		}
	}

	friends := make([]model.User, 0)
	for _, friendID := range friendIDs {
		friend, err := ur.GetByID(friendID)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (ur *UserRepositoryImpl) GetBlockedUsersByID(id string) ([]model.User, error) {
	query := `SELECT user1_id, user2_id FROM block_list WHERE user1_id = ? OR user2_id = ?`
	rows, err := ur.db.Query(query, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockedIDs []string
	for rows.Next() {
		var user1ID, user2ID string
		if err := rows.Scan(&user1ID, &user2ID); err != nil {
			return nil, err
		}
		if user1ID != id {
			blockedIDs = append(blockedIDs, user1ID)
		}
		if user2ID != id {
			blockedIDs = append(blockedIDs, user2ID)
		}
	}

	blocked := make([]model.User, 0)
	for _, blockedID := range blockedIDs {
		blockedUser, err := ur.GetByID(blockedID)
		if err != nil {
			return nil, err
		}
		blocked = append(blocked, blockedUser)
	}

	return blocked, nil
}

func (ur *UserRepositoryImpl) GetByID(id string) (model.User, error) {
	query := `SELECT id, user_id, name FROM users WHERE id = ?`
	row := ur.db.QueryRow(query, id)

	var user model.User
	if err := row.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, fmt.Errorf("user not found")
		}
		return model.User{}, err
	}
	return user, nil
}

/*func (ur *UserRepositoryImpl) GetFriendNames(ids []string) ([]string, error) {
	query := `SELECT name FROM users WHERE id IN (?)`
	query = replacePlaceholders(query, len(ids))

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := ur.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		friendNames = append(friendNames, name)
	}

	return friendNames, nil
}*/

func replacePlaceholders(query string, argCount int) string {
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = "?"
	}
	return strings.Replace(query, "?", strings.Join(placeholders, ","), -1)
}
