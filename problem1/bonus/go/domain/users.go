package domain

type UserCollection struct {
	Users []User `json:"users"`
}

func NewUserCollection(users []User) *UserCollection {
	return &UserCollection{users}
}

func (uc UserCollection) GetUniqueUsers() *UserCollection {
	var uniqueUsers []User
	var uniqueUserIDs []UserID
	for _, u := range uc.Users {
		if !contains(uniqueUserIDs, u.UserID) {
			uniqueUsers = append(uniqueUsers, u)
			uniqueUserIDs = append(uniqueUserIDs, u.UserID)
		}
	}

	return NewUserCollection(uniqueUsers)
}

func (uc *UserCollection) GetUserNames() []string {
	userNames := make([]string, 0)
	for _, u := range uc.Users {
		userNames = append(userNames, u.Name)
	}
	return userNames
}
