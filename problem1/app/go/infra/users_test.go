package infra

import (
	"database/sql"
	"fmt"
	"log"
	"problem1/testutils"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type txAdmin struct {
	*sql.DB
	*testing.T
}

func NewTxAdmin(db *sql.DB, t *testing.T) *txAdmin {
	return &txAdmin{
		db,
		t,
	}
}

// Transaction for test
func (ta *txAdmin) Transaction(update func(tx *sql.Tx) (err error)) error {
	tx, err := ta.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	_ = testutils.PrepareTestBlockList(ta.T, tx)
	_ = testutils.PrepareTestFriendLinks(ta.T, tx)
	_ = testutils.PrepareTestUsers(ta.T, tx)
	ta.T.Cleanup(func() { _ = tx.Rollback() })

	if err := update(tx); err != nil {
		return fmt.Errorf("transaction query failed %w", err)
	}
	// Test用なのでCommitしない
	return nil
}

func TestUserRepository_GetByID(t *testing.T) {
	db := testutils.OpenDBForTest(t)
	ta := NewTxAdmin(db, t)

	sut := NewUserRepository()

	ret, err := sut.GetByID(1, ta)
	if err != nil {
		t.Fatal(err)
	}

	if ret.UserID != 1 {
		t.Errorf("ID should be %d, but got %d", 1, ret.ID)
	}

	log.Printf("ret: %+v", ret)
}

func TestGetFriendsByID(t *testing.T) {
	tx, err := testutils.OpenDBForTest(t).Begin()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	_ = testutils.PrepareTestBlockList(t, tx)
	_ = testutils.PrepareTestFriendLinks(t, tx)
	_ = testutils.PrepareTestUsers(t, tx)

	sut := NewUserRepository()
	ret, err := sut.getFriendsByID(1, tx)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret) != 4 {
		t.Errorf("friends length should be 4, but got %d", len(ret))
	}

	log.Printf("ret: %+v", ret)
}
