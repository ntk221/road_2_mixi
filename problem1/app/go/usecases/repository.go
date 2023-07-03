package usecases

import (
	"problem1/domain"
)

type UserRepository interface {
	UserGetter
}

type UserGetter interface {
	GetByID(id int, db domain.Queryer) (domain.User, error)
}
