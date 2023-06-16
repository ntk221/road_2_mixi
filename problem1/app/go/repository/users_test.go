package repository

import (
	"database/sql"
	"problem1/configs"
	"problem1/model"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestUserRepository_GetByID(t *testing.T) {
	t.Helper()

	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatal(err)
	}

	// 最初にDBをクリアする
	if _, err := db.Exec(`DELETE FROM users;`); err != nil {
		t.Logf("failed to delete user: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = tx.Rollback()
		db.Close()
	})

	// テスト用のデータを作成
	testUser := model.User{ID: 1, UserID: 1, Name: "Test User"}
	_, err = tx.Exec(`
		INSERT INTO users (id, user_id, name) VALUES (?, ?, ?);
	`, testUser.ID, testUser.UserID, testUser.Name)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	row := tx.QueryRow("SELECT COUNT(*) FROM users")
	err = row.Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("unexpected number of rows inserted: got %d, want %d", count, 1)
	}

	sut := NewUserRepository()

	got, err := sut.GetByID(testUser.UserID, tx)

	if err != nil {
		t.Errorf("got %v, want %v", err, testUser)
	}

	if got.ID != testUser.ID {
		t.Errorf("got %v, want %v", err, testUser)
	}

	// 存在しないユーザーの取得をテスト
	/*notFoundUserID := "notfounduser"
	_, err = sut.GetByID(notFoundUserID)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("got %v, want user not found error", err)
	}*/
}

// テスト用のユーザーをDBに登録する
/*func prepareUsers(t *testing.T) []model.User {
	t.Helper()

	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatal(err)
	}

	// 最初にDBをクリアする
	if _, err := db.Exec(`DELETE FROM user`); err != nil {
		t.Logf("failed to delete user: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })

	wants := []model.User{
		{ID: 1, Name: "user1"},
		{ID: 2, Name: "user2"},
		{ID: 3, Name: "user3"},
	}

	_, err = tx.Exec(`
		INSERT INTO user (id, name) VALUES
			(?, ?),
			(?, ?),
			(?, ?);`,
		wants[0].ID, wants[0].Name,
		wants[1].ID, wants[1].Name,
		wants[2].ID, wants[2].Name,
	)

	if err != nil {
		t.Fatal(err)
	}

	return wants
}
*/
