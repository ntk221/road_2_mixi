package repository

import (
	"database/sql"
	"problem1/configs"
	"problem1/model"
	"problem1/types"
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

func TestUserRepository_GetFriendsByID(t *testing.T) {
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

	if _, err := db.Exec(`DELETE FROM friend_link;`); err != nil {
		t.Logf("failed to delete friend_link: %v", err)
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
	testUsers := []model.User{
		{ID: 1, UserID: 1, Name: "Test User1"},
		{ID: 2, UserID: 2, Name: "Test User2"},
		{ID: 3, UserID: 3, Name: "Test User3"},
		{ID: 4, UserID: 4, Name: "Test User4"},
		{ID: 5, UserID: 5, Name: "Test User5"},
		{ID: 6, UserID: 6, Name: "Test User6"},
		{ID: 7, UserID: 7, Name: "Test User7"},
		{ID: 8, UserID: 8, Name: "Test User8"},
		{ID: 9, UserID: 9, Name: "Test User9"},
		{ID: 10, UserID: 10, Name: "Test User10"},
	}
	_, err = tx.Exec(`
	INSERT INTO users (id, user_id, name) VALUES (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?);
`,
		testUsers[0].ID, testUsers[0].UserID, testUsers[0].Name,
		testUsers[1].ID, testUsers[1].UserID, testUsers[1].Name,
		testUsers[2].ID, testUsers[2].UserID, testUsers[2].Name,
		testUsers[3].ID, testUsers[3].UserID, testUsers[3].Name,
		testUsers[4].ID, testUsers[4].UserID, testUsers[4].Name,
		testUsers[5].ID, testUsers[5].UserID, testUsers[5].Name,
		testUsers[6].ID, testUsers[6].UserID, testUsers[6].Name,
		testUsers[7].ID, testUsers[7].UserID, testUsers[7].Name,
		testUsers[8].ID, testUsers[8].UserID, testUsers[8].Name,
		testUsers[9].ID, testUsers[9].UserID, testUsers[9].Name,
	)

	if err != nil {
		t.Fatal(err)
	}

	testFriendLink := []model.FriendLink{
		{ID: 1, User1ID: 1, User2ID: 2},
		{ID: 2, User1ID: 1, User2ID: 3},
		{ID: 3, User1ID: 1, User2ID: 4},
		{ID: 4, User1ID: 1, User2ID: 5},
		{ID: 5, User1ID: 1, User2ID: 6},
		{ID: 6, User1ID: 1, User2ID: 7},
		{ID: 7, User1ID: 1, User2ID: 8},
		{ID: 8, User1ID: 1, User2ID: 9},
		{ID: 9, User1ID: 1, User2ID: 10},
	}

	_, err = tx.Exec(`
	INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?),
	 (?, ?),
	 (?, ?),
	 (?, ?),
	 (?, ?),
	 (?, ?),
	 (?, ?),
	 (?, ?),
	 (?, ?);
`,
		testFriendLink[0].User1ID, testFriendLink[0].User2ID,
		testFriendLink[1].User1ID, testFriendLink[1].User2ID,
		testFriendLink[2].User1ID, testFriendLink[2].User2ID,
		testFriendLink[3].User1ID, testFriendLink[3].User2ID,
		testFriendLink[4].User1ID, testFriendLink[4].User2ID,
		testFriendLink[5].User1ID, testFriendLink[5].User2ID,
		testFriendLink[6].User1ID, testFriendLink[6].User2ID,
		testFriendLink[7].User1ID, testFriendLink[7].User2ID,
		testFriendLink[8].User1ID, testFriendLink[8].User2ID,
	)

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

	testParams := types.PagenationParams{Limit: 2, Offset: 1}

	got, err := sut.GetFriendsByID(testUsers[0].UserID, testParams, tx)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("unexpected number of rows inserted: got %d, want %d", len(got), 2)
	}
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
}*/
