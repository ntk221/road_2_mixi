package repository

import (
	"problem1/domain/entity"
	"problem1/domain/valueObject"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . UserRepository UserGetter
type UserRepository interface {
	UserGetter
}

type UserGetter interface {
	GetByID(id valueObject.UserID, db Queryer) (entity.User, error)
}
