package handler

import (
	"database/sql"
	"fmt"
)

type txAdmin struct {
	*sql.DB
}

func NewTxAdmin(db *sql.DB) *txAdmin {
	return &txAdmin{db}
}

func (ta *txAdmin) Transaction(update func() (err error)) error {
	tx, err := ta.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	if err := update(); err != nil {
		return fmt.Errorf("transaction query failed %w", err)
	}
	return tx.Commit()
}
