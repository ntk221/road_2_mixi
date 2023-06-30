package domain

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

func (u User) GetFriendList() []int {
	return u.FriendList
}

func (u User) GetBlockList() []int {
	return u.FriendList
}
