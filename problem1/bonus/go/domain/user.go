package domain

type UserID int

type User struct {
	ID         int64  `db:"id"`
	UserID     UserID `db:"user_id"`
	Name       string `db:"name"`
	FriendList []UserID
	BlockList  []UserID
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

func (u User) GetFriendList() []UserID {
	return u.FriendList
}

func (u User) GetBlockList() []UserID {
	return u.BlockList
}

func (u User) IsBlocked(user User) bool {
	return contains(user.BlockList, u.UserID)
}

func contains(list []UserID, target UserID) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
