package infra

import (
	"database/sql"
	"errors"
	"fmt"
	"problem1/domain"
)

type UserRepositoryImpl struct {
	domain.UserGetter
}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

// FriendLinkによって，繋がっているユーザーのIDを取得する
// 方向性は考慮しない
func (ur *UserRepositoryImpl) getFriendsByID(user_id domain.UserID, db domain.Queryer) ([]domain.UserID, error) {
	query := `
		SELECT user_id
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
	`

	// params.Limit = 2

	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	friends := make([]domain.UserID, 0)
	for rows.Next() {
		var friend domain.User
		if err := rows.Scan(&friend.UserID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, sql.ErrNoRows
			}
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		friends = append(friends, friend.UserID)
	}

	return friends, nil
}

// BlockListによって繋がっている，ユーザーのIDを取得する
// 方向性は考慮しない
func (ur *UserRepositoryImpl) getBlockUsersByID(user_id domain.UserID, db domain.Queryer) ([]domain.UserID, error) {
	query := `
		SELECT user1_id, user2_id 
		FROM block_list 
		WHERE user1_id = ? 
		OR user2_id = ?
	`
	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var blockIDs []domain.UserID
	for rows.Next() {
		var user1ID, user2ID domain.UserID
		if err := rows.Scan(&user1ID, &user2ID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, sql.ErrNoRows
			}
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if user1ID != user_id {
			blockIDs = append(blockIDs, user1ID)
		}
		if user2ID != user_id {
			blockIDs = append(blockIDs, user2ID)
		}
	}

	return blockIDs, nil
}

func (ur *UserRepositoryImpl) GetByID(user_id domain.UserID, db domain.Queryer) (domain.User, error) {
	var user domain.User

	query := `SELECT id, user_id, name FROM users WHERE user_id = ?`
	row := db.QueryRow(query, user_id)

	if err := row.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
		return domain.User{}, fmt.Errorf("failed to scan row: %w", err)
	}

	friendIDs, err := ur.getFriendsByID(user_id, db)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get friends: %w", err)
	}
	user.FriendList = friendIDs

	blockedIDs, err := ur.getBlockUsersByID(user_id, db)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get blocked users: %w", err)
	}
	user.BlockList = blockedIDs

	user = *domain.NewUser(user.ID, user.UserID, user.Name, user.FriendList, user.BlockList)

	return user, nil
}

func (ur *UserRepositoryImpl) GetUsers(db domain.Queryer) ([]domain.User, error) {
	users := make([]domain.User, 0)
	query := `
		SELECT id, user_id, name
		FROM users
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		friendIDs, err := ur.getFriendsByID(user.UserID, db)
		if err != nil {
			return nil, fmt.Errorf("failed to get friends: %w", err)
		}
		user.FriendList = friendIDs

		blockedIDs, err := ur.getBlockUsersByID(user.UserID, db)
		if err != nil {
			return nil, fmt.Errorf("failed to get blocked users: %w", err)
		}
		user.BlockList = blockedIDs
		user = *domain.NewUser(user.ID, user.UserID, user.Name, user.FriendList, user.BlockList)
		users = append(users, user)
	}

	return users, nil
}
