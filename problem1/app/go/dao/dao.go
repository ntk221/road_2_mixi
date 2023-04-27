package dao

import (
	"database/sql"
	"problem1/domain/repository"

	_ "github.com/go-sql-driver/mysql"
)

type DAO interface {
	User() repository.User
}

type dao struct {
	db *sql.DB
}

func New(db *sql.DB) DAO {
	return &dao{db: db}
}

func (d *dao) User() repository.User {
	return NewUser(d.db)
}
