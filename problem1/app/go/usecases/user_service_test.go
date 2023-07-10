package usecases_test

import (
	"database/sql"
	"fmt"
	"log"
	"problem1/infra"
	"problem1/testutils"
	"problem1/usecases"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type txAdmin struct {
	*sql.DB
	*sql.Tx
	*testing.T
}

func NewTxAdmin(db *sql.DB, t *testing.T) *txAdmin {
	return &txAdmin{
		db,
		nil,
		t,
	}
}

func (ta *txAdmin) Prepare(query string) (*sql.Stmt, error) {
	return ta.Tx.Prepare(query)
}

func (ta *txAdmin) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return ta.Tx.Query(query, args...)
}

func (ta *txAdmin) QueryRow(query string, args ...interface{}) *sql.Row {
	return ta.Tx.QueryRow(query, args...)
}

// Transaction for test
func (ta *txAdmin) Transaction(update func(tx *sql.Tx) (err error)) error {
	if ta.Tx == nil {
		tx, err := ta.DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		ta.Tx = tx
		_ = testutils.PrepareTestBlockList(ta.T, tx)
		_ = testutils.PrepareTestFriendLinks(ta.T, tx)
		_ = testutils.PrepareTestUsers(ta.T, tx)
		ta.T.Cleanup(func() { _ = ta.Tx.Rollback() })
	}

	if err := update(ta.Tx); err != nil {
		return fmt.Errorf("transaction query failed %w", err)
	}
	// Test用なのでCommitしない
	return nil
}

func TestGetRealFriends(t *testing.T) {
	db := testutils.OpenDBForTest(t)

	ta := NewTxAdmin(db, t)
	ur := infra.NewUserRepository()
	sut := usecases.NewUserService(ta, ur)

	ret, err := sut.GetFriendList(1)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%v\n", ret)

}
