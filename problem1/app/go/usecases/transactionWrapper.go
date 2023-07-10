package usecases

import (
	"database/sql"
	"fmt"
	"problem1/domain"
)

type txAdmin struct {
	domain.Database
}

func NewTxAdmin(db domain.Database) *txAdmin {
	return &txAdmin{db}
}

func (ta *txAdmin) Transaction(update func(tx *sql.Tx) (err error)) error {
	tx, err := ta.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	if err := update(tx); err != nil {
		return fmt.Errorf("transaction query failed %w", err)
	}
	return tx.Commit()
}
