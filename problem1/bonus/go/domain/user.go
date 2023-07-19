package domain

type UserID int

type User struct {
	ID         int64
	UserID     UserID
	Name       string
	FriendList []UserID
	BlockList  []UserID
}

func NewUser(ID int64, UserID UserID, Name string, FriendList []UserID, BlockList []UserID) *User {
	u := &User{
		ID,
		UserID,
		Name,
		FriendList,
		BlockList,
	}

	// ブロックリストに入っているユーザーは，友達リストから除外する
	u = u.filterByBlockList()
	return u
}

func (u User) filterByBlockList() *User {
	var filteredFriendList []UserID
	for _, v := range u.FriendList {
		if !contains(u.BlockList, v) {
			filteredFriendList = append(filteredFriendList, v)
		}
	}
	u.FriendList = filteredFriendList
	return &u
}

func (u User) GetFriendList() []UserID {
	return u.FriendList
}

func (u User) GetBlockList() []UserID {
	return u.BlockList
}

/*func (u User) IsBlocked(user User) bool {
	return contains(user.BlockList, u.UserID)
}*/

func contains[T comparable](list []T, target T) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
