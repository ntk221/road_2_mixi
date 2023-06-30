package domain

import "fmt"

type User struct {
	ID         int64  `db:"id"`
	UserID     int    `db:"user_id"`
	Name       string `db:"name"`
	FriendList []int
	BlockList  []int
}

type FriendLink struct {
	ID      int64 `db:"id"`
	User1ID int   `db:"user1_id"`
	User2ID int   `db:"user2_id"`
}

type BlockList struct {
	ID      int64 `db:"id"`
	User1ID int   `db:"user1_id"`
	User2ID int   `db:"user2_id"`
}

func (u User) GetFriendList() ([]int, error) {
	if u.FriendList == nil {
		return nil, fmt.Errorf("%v has no friend", u.Name)
	}
	return u.FriendList, nil
}

func (u User) GetBlockList() ([]int, error) {
	if u.FriendList == nil {
		return nil, fmt.Errorf("%v block no one", u.Name)
	}
	return u.FriendList, nil
}
