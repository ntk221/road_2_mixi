package infra

import (
	"database/sql"
	"errors"
	"fmt"
	"problem1/domain"
)

type UserRepository interface {
	GetFriendsByID(user_id int, db Queryer) ([]domain.User, error)
	GetBlockUsersByID(user_id int, db Queryer) ([]domain.User, error)
	GetByID(user_id int, db Queryer) (domain.User, error)
	// GetFriendNames(ids []string) ([]string, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

// FriendLinkテーブルから友達のIDを取得
// 一方が友達になっている場合には，友人リストに追加している
func (ur *UserRepositoryImpl) GetFriendsByID(user_id int, db Queryer) ([]int, error) {
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

	friends := make([]int, 0)
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

// ブロックしているユーザーのIDを取得
func (ur *UserRepositoryImpl) GetBlockUsersByID(user_id int, db Queryer) ([]int, error) {
	query := `SELECT user1_id, user2_id FROM block_list WHERE user1_id = ? OR user2_id = ?`
	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var blockIDs []int
	for rows.Next() {
		var user1ID, user2ID int
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

func (ur *UserRepositoryImpl) GetByID(user_id int, db Queryer) (domain.User, error) {
	query := `SELECT id, user_id, name FROM users WHERE user_id = ?`
	row := db.QueryRow(query, user_id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 存在しないユーザーIDが指定された場合はpanic
			// これの代わりにエラーを返した方が良いかも
			return domain.User{}, sql.ErrNoRows
		}
		return domain.User{}, fmt.Errorf("failed to scan row: %w", err)
	}

	friendIDs, err := ur.GetFriendsByID(user_id, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user.FriendList = nil // return model.User{}, sql.ErrNoRows
		}
		return domain.User{}, fmt.Errorf("failed to get friends: %w", err)
	}
	user.FriendList = friendIDs

	blockedIDs, err := ur.GetBlockUsersByID(user_id, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user.BlockList = nil // return model.User{}, sql.ErrNoRows
		}
		return domain.User{}, fmt.Errorf("failed to get blocked users: %w", err)
	}
	user.BlockList = blockedIDs

	return user, nil
}

/*func replacePlaceholders(query string, argCount int) string {
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = "?"
	}
	return strings.Replace(query, "?", strings.Join(placeholders, ","), -1)
}*/
