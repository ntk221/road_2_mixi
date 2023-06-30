package usecases

import (
	"problem1/domain"
	"problem1/infra"
)

type UserGetter interface {
	GetByID(id int, db infra.Queryer) (domain.User, error)
}
