package infra

import (
	"database/sql"
	"fmt"
	"log"
	"problem1/domain"
	"problem1/testutils"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type txAdmin struct {
	*sql.DB
}

func NewTxAdmin(db *sql.DB) *txAdmin {
	return &txAdmin{db}
}

// Transaction for test
func (ta *txAdmin) Transaction(update func() (err error)) error {
	tx, err := ta.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() 
	if err := update(); err != nil {
		return fmt.Errorf("transaction query failed %w", err)
	}
	// Test用なのでCommitしない
	return nil
}

func TestUserRepository_GetByID(t *testing.T) {
	db := testutils.OpenDBForTest(t)
	ta := NewTxAdmin(db)

	testUsers := []domain.User{
		{
			UserID: 1,
			Name:   "user1",
			FriendList: []int{
				2, 3, 4, 5,
			},
			BlockList: []int{
				6, 7, 8, 9,
			},
		},
		{
			UserID: 2,
			Name:   "user2",
			FriendList: []int{
				1, 3, 4, 5,
			},
			BlockList: []int{
				6, 7, 8, 9,
			},
		},
	}

	log.Printf("testUsers: %+v", testUsers)

	sut := NewUserRepository()

	ret, err := sut.GetByID(testUsers[0].UserID, ta)
	if err != nil {
		t.Fatal(err)
	}

	if ret.UserID != testUsers[0].UserID {
		t.Errorf("ID should be %d, but got %d", testUsers[0].ID, ret.ID)
	}

	log.Printf("ret: %+v", ret)
}

func TestGetFriendsByID(t *testing.T) {
	tx, err := testutils.OpenDBForTest(t).Begin()
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	testUsers := testutils.PrepareTestUsers(t, tx)
	testFriendLink := testutils.PrepareTestFriendLinks(t, tx)
	_ = testFriendLink

	sut := NewUserRepository()
	ret, err := sut.getFriendsByID(testUsers[0].UserID, tx)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret) != 4 {
		t.Errorf("friends length should be 4, but got %d", len(ret))
	}
}
