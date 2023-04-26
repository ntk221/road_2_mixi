package main

import (
	"database/sql"
	"net/http"
	"problem1/configs"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "minimal_sns_app")
	})

	e.GET("/get_friend_list", func(c echo.Context) error {
		id := c.QueryParam("id")
		friendIdList, err := getFriendIdList(db, id)
		if err != nil {
			return err
		}
		friendList, err := getFriendList(db, friendIdList)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, friendList)
	})

	e.GET("/get_friend_of_friend_list", func(c echo.Context) error {
		id := c.QueryParam("id")
		friendIdList, err := getFriendIdList(db, id)
		if err != nil {
			return err
		}
		friendsOfFriends, err := getFriendsOfFriends(db, friendIdList)
		if err != nil {
			return err
		}

		var ret []string
		for friend := range friendsOfFriends {
			ret = append(ret, friend)
		}
		return c.JSON(http.StatusOK, ret)
	})

	e.GET("/get_friend_of_friend_list_paging", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}

// 友達IDリストから友達リストを取得する
func getFriendList(db *sql.DB, friendIdList []string) ([]string, error) {
	var friendList []string
	for _, friendID := range friendIdList {
		rows, err := db.Query("SELECT name FROM users WHERE user_id = ?", friendID)
		if err != nil {
			return nil, err
		}
		var friend string
		rows.Next()
		if err := rows.Scan(&friend); err != nil {
			return nil, err
		}
		friendList = append(friendList, friend)
	}
	return friendList, nil
}

// 友達リストからリストの各要素について，その友達リストを取得する
func getFriendsOfFriends(db *sql.DB, friendIdList []string) (map[string]bool, error) {
	friendsOfFriends := make(map[string]bool)
	for _, friendId := range friendIdList {
		friendsId, err := getFriendIdList(db, friendId)
		if err != nil {
			return nil, err
		}
		friends, err := getFriendList(db, friendsId)
		if err != nil {
			return nil, err
		}
		for _, friend := range friends {
			if _, ok := friendsOfFriends[friend]; !ok {
				friendsOfFriends[friend] = true
			}
		}
	}
	return friendsOfFriends, nil
}

func getFriendIdList(db *sql.DB, id string) ([]string, error) {
	query := `
		SELECT user2_id FROM friend_link WHERE user1_id = ?
		UNION
		SELECT user1_id FROM friend_link WHERE user2_id = ?
	`
	rows, err := db.Query(query, id, id)
	if err != nil {
		panic(err)
	}
	var friendIdList []string
	for rows.Next() {
		var friendId string
		if err := rows.Scan(&friendId); err != nil {
			return nil, err
		}
		friendIdList = append(friendIdList, friendId)
	}
	return friendIdList, nil
}
