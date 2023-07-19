package repository

import (
	"database/sql"
)

type Database interface {
	Beginner
	Preparer
	Queryer
	Executer
}

type Beginner interface {
	Begin() (*sql.Tx, error)
}

type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

type Executer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Queryer interface {
	Preparer
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

var (
	_ Beginner = (*sql.DB)(nil)
	_ Preparer = (*sql.DB)(nil)
	_ Queryer  = (*sql.DB)(nil)
	_ Queryer  = (*sql.Tx)(nil)
	_ Executer = (*sql.DB)(nil)
	_ Executer = (*sql.Tx)(nil)
)
