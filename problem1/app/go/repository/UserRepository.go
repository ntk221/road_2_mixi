package repository

import (
	"database/sql"
	"strings"
)

type UserRepository interface {
	GetFriendIDs(id string) ([]string, error)
	GetBlockedUsers(id string) ([]string, error)
	GetFriendNames(ids []string) ([]string, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) GetFriendIDs(id string) ([]string, error) {
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

	return friendIDs, nil
}

func (ur *UserRepositoryImpl) GetBlockedUsers(id string) ([]string, error) {
	query := `SELECT user1_id, user2_id FROM block_list WHERE user1_id = ? OR user2_id = ?`
	rows, err := ur.db.Query(query, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockedUsers []string
	for rows.Next() {
		var user1ID, user2ID string
		if err := rows.Scan(&user1ID, &user2ID); err != nil {
			return nil, err
		}
		if user1ID != id {
			blockedUsers = append(blockedUsers, user1ID)
		}
		if user2ID != id {
			blockedUsers = append(blockedUsers, user2ID)
		}
	}

	return blockedUsers, nil
}

func (ur *UserRepositoryImpl) GetFriendNames(ids []string) ([]string, error) {
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
}

func replacePlaceholders(query string, argCount int) string {
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = "?"
	}
	return strings.Replace(query, "?", strings.Join(placeholders, ","), -1)
}
