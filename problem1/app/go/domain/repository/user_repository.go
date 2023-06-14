package repository

import (
	"problem1/domain/object"
)

type User interface {
	GetByID(id string) (*object.User, error)
	// GetByName(name string) (*object.User, error)
	GetFriends(id string) ([]*object.User, error)
}

// DIPを採用したことによって，daoへの依存は circuler dependency になる
/*
func NewUserRepository(db *sql.DB) User {
	d := dao.New(db)
	return d.User()
}
*/
