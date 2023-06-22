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
	testUser := model.User{ID: 1, UserID: 1, Name: "Test User", FriendList: []int{2, 3}, BlockList: []int{4, 5}}
	_, err = tx.Exec(`
		INSERT INTO users (id, user_id, name) VALUES (?, ?, ?);
	`, testUser.ID, testUser.UserID, testUser.Name)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec(`
		INSERT INTO friend_links (user1_id, user2_id) VALUES (?, ?), (?, ?);
	`, testUser.UserID, testUser.FriendList[0], testUser.UserID, testUser.FriendList[1])
	if err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec(`
		INSERT INTO block_list (user1_id, user2_id) VALUES (?, ?), (?, ?);
	`, testUser.UserID, testUser.BlockList[0], testUser.UserID, testUser.BlockList[1])
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
	// 存在しないユーザーに関して即座にpanicはどうなんだろう
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
	}

	_, err = tx.Exec(`
		INSERT INTO users (id, user_id, name) VALUES (?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?)
	;
	`, testUsers[0].ID, testUsers[0].UserID, testUsers[0].Name,
		testUsers[1].ID, testUsers[1].UserID, testUsers[1].Name,
		testUsers[2].ID, testUsers[2].UserID, testUsers[2].Name,
		testUsers[3].ID, testUsers[3].UserID, testUsers[3].Name,
		testUsers[4].ID, testUsers[4].UserID, testUsers[4].Name,
	)

	if err != nil {
		t.Fatal(err)
	}

	testLinks := []model.FriendLink{
		{ID: 1, User1ID: 1, User2ID: 2},
		{ID: 2, User1ID: 1, User2ID: 3},
		{ID: 3, User1ID: 2, User2ID: 1},
		{ID: 4, User1ID: 3, User2ID: 2},
		{ID: 5, User1ID: 3, User2ID: 1},
		{ID: 6, User1ID: 4, User2ID: 2},
		{ID: 7, User1ID: 4, User2ID: 1},
		{ID: 8, User1ID: 5, User2ID: 2},
		{ID: 9, User1ID: 5, User2ID: 1},
	}

	_, err = tx.Exec(`
		INSERT INTO friend_link (id, user1_id, user2_id) VALUES (?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?)
	`, testLinks[0].ID, testLinks[0].User1ID, testLinks[0].User2ID,
		testLinks[1].ID, testLinks[1].User1ID, testLinks[1].User2ID,
		testLinks[2].ID, testLinks[2].User1ID, testLinks[2].User2ID,
		testLinks[3].ID, testLinks[3].User1ID, testLinks[3].User2ID,
		testLinks[4].ID, testLinks[4].User1ID, testLinks[4].User2ID,
		testLinks[5].ID, testLinks[5].User1ID, testLinks[5].User2ID,
		testLinks[6].ID, testLinks[6].User1ID, testLinks[6].User2ID,
		testLinks[7].ID, testLinks[7].User1ID, testLinks[7].User2ID,
		testLinks[8].ID, testLinks[8].User1ID, testLinks[8].User2ID,
	)

	if err != nil {
		t.Fatal(err)
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
