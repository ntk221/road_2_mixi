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
		friendList, err := getFriendList(db, id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, friendList)
	})

	e.GET("/get_friend_of_friend_list", func(c echo.Context) error {

		return nil
	})

	e.GET("/get_friend_of_friend_list_paging", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}

func getFriendList(db *sql.DB, id string) []string, error {
	rows, err := db.Query("SELECT user2_id FROM friend_link WHERE user1_id = ?", id)
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
	var friendList []string
	for _, friendID := range friendIdList {
		rows, err := db.Query("SELECT name FROM users WHERE user_id = ?", friendID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var friend string
			if err := rows.Scan(&friend); err != nil {
				return nil, err
			}
			friendList = append(friendList, friend)
		}
	}
	return friendList, nil
}
