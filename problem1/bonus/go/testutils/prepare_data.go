package testutils

import (
	"problem1/domain"
	"testing"
)

func PrepareTestUsers(t *testing.T, tx domain.Executer) []domain.User {
	t.Helper()

	if _, err := tx.Exec(`DELETE FROM users`); err != nil {
		t.Fatal(err)
	}

	testUsers := []domain.User{
		{UserID: 1, Name: "Test User1"},
		{UserID: 2, Name: "Test User2"},
		{UserID: 3, Name: "Test User3"},
		{UserID: 4, Name: "Test User4"},
		{UserID: 5, Name: "Test User5"},
	}

	insertUserQuery := `INSERT INTO users (user_id, name) VALUES (?, ?)`
	for _, user := range testUsers {
		_, err := tx.Exec(insertUserQuery, user.UserID, user.Name)
		if err != nil {
			t.Fatal(err)
		}
	}

	return testUsers
}

func PrepareTestFriendLinks(t *testing.T, tx domain.Executer) []domain.FriendLink {
	t.Helper()

	if _, err := tx.Exec(`DELETE FROM friend_link`); err != nil {
		t.Fatal(err)
	}

	testFriendLinks := []domain.FriendLink{
		{User1ID: 1, User2ID: 2},
		{User1ID: 1, User2ID: 3},
		{User1ID: 1, User2ID: 4},
		{User1ID: 2, User2ID: 3},
		{User1ID: 2, User2ID: 4},
		{User1ID: 3, User2ID: 4},
		{User1ID: 5, User2ID: 1},
	}

	insertFriendLinkQuery := `INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?)`
	for _, friendLink := range testFriendLinks {
		_, err := tx.Exec(insertFriendLinkQuery, friendLink.User1ID, friendLink.User2ID)
		if err != nil {
			t.Fatal(err)
		}
	}

	return testFriendLinks
}

func PrepareTestBlockList(t *testing.T, tx domain.Executer) []domain.BlockList {
	t.Helper()

	if _, err := tx.Exec(`DELETE FROM block_list`); err != nil {
		t.Fatal(err)
	}

	testBlockList := []domain.BlockList{
		{User1ID: 1, User2ID: 2},
		{User1ID: 2, User2ID: 3},
		{User1ID: 3, User2ID: 4},
		{User1ID: 4, User2ID: 1},
	}

	insertBlockListQuery := `INSERT INTO block_list (user1_id, user2_id) VALUES (?, ?)`
	for _, blockList := range testBlockList {
		_, err := tx.Exec(insertBlockListQuery, blockList.User1ID, blockList.User2ID)
		if err != nil {
			t.Fatal(err)
		}
	}

	return testBlockList
}
