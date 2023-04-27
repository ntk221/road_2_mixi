package object

type User struct {
	Id      string
	Name    string
	Friends []*User
}

func (u *User) AddFriend(friend *User) {
	u.Friends = append(u.Friends, friend)
}

func (u *User) GetFriends() []*User {
	return u.Friends
}

func (u *User) IsFriend(friend *User) bool {
	for _, f := range u.Friends {
		if f.Id == friend.Id {
			return true
		}
	}
	return false
}
