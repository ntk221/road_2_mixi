package dao

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Friend struct {
	Id   string
	Name string
}

type FriendDAO interface {
	GetFriends(friendIdList []string) (map[Friend]bool, error)
	GetFriendIdList(id string) ([]string, error)
	GetBlockedUsers(id string) (map[string]bool, error)
	GetFriendsOfFriends(friendIdList []string, blockedUsers map[string]bool) (map[Friend]bool, error)
	NewSQLFriendDAO(id string, name string) error
}

type SQLFriendDAO struct {
	DB *sql.DB
}

func NewSQLFriendDAO(db *sql.DB) *SQLFriendDAO {
	return &SQLFriendDAO{DB: db}
}

func (dao *SQLFriendDAO) GetFriends(friendIdList []string) (map[Friend]bool, error) {
	friendList := make(map[Friend]bool, 0)
	for _, friendId := range friendIdList {
		var friend Friend
		row := dao.DB.QueryRow("SELECT user_id, name FROM users WHERE user_id = ?", friendId)
		if err := row.Scan(&friend.Id, &friend.Name); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		friendList[friend] = true
	}
	return friendList, nil
}

// GetBlockedUsers メソッドの実装
func (dao *SQLFriendDAO) GetBlockedUsers(id string) (map[string]bool, error) {
	blockedIds := make(map[string]bool)
	rows, err := dao.DB.Query("SELECT user2_id FROM block_list WHERE user1_id = ? UNION SELECT user1_id FROM block_list WHERE user2_id = ?", id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var blockedId string
		if err := rows.Scan(&blockedId); err != nil {
			return nil, err
		}
		blockedIds[blockedId] = true
	}
	return blockedIds, nil
}

func (dao *SQLFriendDAO) GetFriendIdList(id string) ([]string, error) {
	friendIdList := make([]string, 0)
	blockedList, err := dao.GetBlockedUsers(id)
	if err != nil {
		return nil, err
	}
	rows, err := dao.DB.Query("SELECT user2_id FROM friend_link WHERE user1_id = ? UNION SELECT user1_id FROM friend_link WHERE user2_id = ?", id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var friendId string
		if err := rows.Scan(&friendId); err != nil {
			return nil, err
		}
		if _, yes := blockedList[friendId]; yes {
			continue
		}
		friendIdList = append(friendIdList, friendId)
	}
	return friendIdList, nil
}

func (dao *SQLFriendDAO) GetFriendsOfFriends(friendIdList []string, blockedUsers map[string]bool) (map[Friend]bool, error) {
	friendList := make(map[Friend]bool, 0)
	for _, friendId := range friendIdList {
		friendIdList, err := dao.GetFriendIdList(friendId)
		if err != nil {
			return nil, err
		}
		friends, err := dao.GetFriends(friendIdList)
		if err != nil {
			return nil, err
		}
		for friend, _ := range friends {
			if _, yes := blockedUsers[friend.Id]; yes {
				continue
			}
			friendList[friend] = true
		}
	}
	return friendList, nil
}
