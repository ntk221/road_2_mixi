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

func NewUser(ID int64, UserID UserID, Name string, FriendList []UserID, BlockList []UserID) *User {
	u := User{
		ID,
		UserID,
		Name,
		FriendList,
		BlockList,
	}

	u = u.filterByBlockList()
	return &u
}

func (u User) filterByBlockList() User {
	var filteredFriendList []UserID
	for _, v := range u.FriendList {
		if !contains(u.BlockList, v) {
			filteredFriendList = append(filteredFriendList, v)
		}
	}
	u.FriendList = filteredFriendList
	return u
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
