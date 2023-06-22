package repository

import (
	"database/sql"
	"fmt"
	"problem1/model"
	"problem1/types"
)

type UserRepository interface {
	GetFriendsByID(user_id int, params types.PagenationParams, db Queryer) ([]model.User, error)
	GetBlockedUsersByID(user_id int, db Queryer) ([]model.User, error)
	GetByID(user_id int, db Queryer) (model.User, error)
	// GetFriendNames(ids []string) ([]string, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (ur *UserRepositoryImpl) GetFriendsByID(user_id int, params types.PagenationParams, db Queryer) ([]model.User, error) {
	query := `
		SELECT user_id, name
		FROM users
		WHERE user_id IN (
			SELECT user1_id
			FROM friend_link
			WHERE user2_id = ? 
			UNION 
			SELECT user2_id
			FROM friend_link
			WHERE user1_id = ?
		)
		LIMIT ? OFFSET ?
	`

	// params.Limit = 2

	rows, err := db.Query(query, user_id, user_id, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := make([]model.User, 0)
	for rows.Next() {
		var friend model.User
		if err := rows.Scan(&friend.UserID, &friend.Name); err != nil {
			panic(err)
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (ur *UserRepositoryImpl) GetBlockedUsersByID(user_id int, db Queryer) ([]model.User, error) {
	query := `SELECT user1_id, user2_id FROM block_list WHERE user1_id = ? OR user2_id = ?`
	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var blockedIDs []int
	for rows.Next() {
		var user1ID, user2ID int
		if err := rows.Scan(&user1ID, &user2ID); err != nil {
			panic(err)
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
			panic(err)
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
			panic(fmt.Sprintf("user_id %d not found", user_id))
		}
		panic(err)
	}
	return user, nil
}

/*func replacePlaceholders(query string, argCount int) string {
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = "?"
	}
	return strings.Replace(query, "?", strings.Join(placeholders, ","), -1)
}*/
