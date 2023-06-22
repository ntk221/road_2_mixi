package model

type User struct {
	ID         int64  `db:"id"`
	UserID     int    `db:"user_id"`
	Name       string `db:"name"`
	FriendList []int
	BlockList  []int
}

func (u *User) IsFriend(other *User) bool {
	if contains(u.FriendList, other.UserID) && contains(other.FriendList, u.UserID) {
		return true
	}
	return false
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
