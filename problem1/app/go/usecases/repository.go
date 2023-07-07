package usecases

import (
	"problem1/domain"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . UserRepository UserGetter
type UserRepository interface {
	UserGetter
}

type UserGetter interface {
	GetByID(id int, db domain.Queryer) (domain.User, error)
}
