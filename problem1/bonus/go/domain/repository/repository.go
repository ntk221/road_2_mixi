package repository

import (
	"problem1/domain"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . UserRepository UserGetter
type UserRepository interface {
	UserGetter
}

type UserGetter interface {
	GetByID(id domain.UserID, db Queryer) (domain.User, error)
}
