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
	return u.BlockList
}

func (u User) IsBlocked(user User) bool {
	return contains(user.BlockList, u.UserID)
}

func contains(list []int, target int) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
