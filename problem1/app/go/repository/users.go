package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"problem1/model"
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

	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		panic(err)
	}
	defer rows.Close()

	friends := make([]model.User, 0)
	for rows.Next() {
		var friend model.User
		if err := rows.Scan(&friend.UserID, &friend.Name); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, sql.ErrNoRows
			}
			return nil, fmt.Errorf("failed to scan row: %w", err)
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
			if errors.Is(err, sql.ErrNoRows) {
				return nil, sql.ErrNoRows
			}
			return nil, fmt.Errorf("failed to scan row: %w", err)
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
			if errors.Is(err, sql.ErrNoRows) {
				return nil, sql.ErrNoRows
			} else {
				return nil, fmt.Errorf("failed to get blocked user: %w", err)
			}
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
		if errors.Is(err, sql.ErrNoRows) {
			// 存在しないユーザーIDが指定された場合はpanic
			// これの代わりにエラーを返した方が良いかも
			return model.User{}, sql.ErrNoRows
		}
		return model.User{}, fmt.Errorf("failed to scan row: %w", err)
	}

	var friends []model.User
	friends, err := ur.GetFriendsByID(user_id, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user_id %d is not found", user_id)
		}
		return model.User{}, fmt.Errorf("failed to get friends: %w", err)
	}
	for _, friend := range friends {
		user.FriendList = append(user.FriendList, friend.UserID)
	}

	var blocked []model.User
	blocked, err = ur.GetBlockedUsersByID(user_id, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user_id %d is not found", user_id)
		}
		return model.User{}, fmt.Errorf("failed to get blocked users: %w", err)
	}

	for _, blockedUser := range blocked {
		user.BlockList = append(user.BlockList, blockedUser.UserID)
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
