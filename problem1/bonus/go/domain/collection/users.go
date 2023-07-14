package domain

import (
	"problem1/domain"
)

type UserCollection struct {
	Users []domain.User `json:"users"`
}

/*func NewUsers(userIDs []domain.UserID) *UserCollection {
}*/
