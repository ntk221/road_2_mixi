package dao

import (
	"database/sql"
	"problem1/domain/object"
	"problem1/domain/repository"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	db *sql.DB
}

func NewUser(db *sql.DB) repository.User {
	return &user{db: db}
}

func (r *user) GetByID(id string) (*object.User, error) {
	row := r.db.QueryRow("SELECT user_id, name FROM users WHERE user_id = ?", id)
	var user object.User
	if err := row.Scan(&user.Id, &user.Name); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *user) GetFriends(id string) ([]*object.User, error) {
	query := `
		SELECT user1_id FROM friend_link WHERE user2_id = ?
		UNION
		SELECT user2_id FROM friend_link WHERE user1_id = ?
	`
	rows, err := r.db.Query(query, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendsId []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		friendsId = append(friendsId, id)
	}

	var friends []*object.User
	for _, friendId := range friendsId {
		friend, err := r.GetByID(friendId)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}
