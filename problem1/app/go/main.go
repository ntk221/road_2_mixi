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
		} else if friendIdList == nil {
			return c.JSON(http.StatusOK, "no friends")
		}
		friendList, err := getFriendList(db, friendIdList)
		if err != nil {
			return err
		} else if friendList == nil {
			return c.JSON(http.StatusOK, "no friends")
		}
		ret := make([]string, 0, len(friendList))
		for friend := range friendList {
			ret = append(ret, friend)
		}
		return c.JSON(http.StatusOK, ret)
	})

	e.GET("/get_friend_of_friend_list", func(c echo.Context) error {
		id := c.QueryParam("id")
		friendIdList, err := getFriendIdList(db, id)
		if err != nil {
			return err
		} else if friendIdList == nil {
			return c.JSON(http.StatusOK, "no friends")
		}
		friendsOfFriends, err := getFriendsOfFriends(db, friendIdList)
		if err != nil {
			return err
		} else if friendsOfFriends == nil {
			return c.JSON(http.StatusOK, "no friends of friends")
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

// 現在は，ビジネスロジック層とデータアクセス層が混在しているが，
// 本来は，データアクセス層はビジネスロジック層に対して，データを返すだけの役割を持つべきである．

// 友達リストからリストの各要素について，その友達リストを取得する
func getFriendsOfFriends(db *sql.DB, friendIdList []string) (map[string]bool, error) {
	if len(friendIdList) == 0 {
		return nil, nil
	}
	friendsOfFriends := make(map[string]bool)
	friendList, err := getFriendList(db, friendIdList)
	if err != nil {
		return nil, err
	}
	for _, friendId := range friendIdList {
		idList, err := getFriendIdList(db, friendId)
		if err != nil {
			return nil, err
		}
		l, err := getFriendList(db, idList)
		if err != nil {
			return nil, err
		}
		for k, _ := range l {
			// 1hop の友達は除外
			if _, ok := friendList[k]; ok {
				continue
			}
			if _, ok := friendsOfFriends[k]; !ok {
				friendsOfFriends[k] = true
			}
		}
	}
	return friendsOfFriends, nil
}

// (ブロック関係を除いた)友達のリストを取得する
func getFriendList(db *sql.DB, friendIdList []string) (map[string]bool, error) {
	friendList := make(map[string]bool)
	for _, friendId := range friendIdList {
		var friend string
		rows, err := db.Query("SELECT name FROM users WHERE user_id = ?", friendId)
		if err != nil {
			return nil, err
		}
		rows.Next()
		if err := rows.Scan(&friend); err != nil {
			return nil, err
		}
		friendList[friend] = true
	}

	return friendList, nil
}

// (ブロック関係を除いた)友達のIDリストを取得する
func getFriendIdList(db *sql.DB, id string) ([]string, error) {
	query := `
		SELECT user2_id FROM friend_link WHERE user1_id = ?
		UNION
		SELECT user1_id FROM friend_link WHERE user2_id = ?
	`
	blockedIds, err := getBlockedIdList(db, id)
	if err != nil {
		return nil, err
	} else if blockedIds == nil {
		return nil, nil
	}
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
		// ブロックされているユーザは除外
		if _, yes := blockedIds[friendId]; yes {
			continue
		}
		friendIdList = append(friendIdList, friendId)
	}
	return friendIdList, nil
}

// ブロック関係にあるユーザのIDリストを取得する
func getBlockedIdList(db *sql.DB, id string) (map[string]bool, error) {
	query := `
		SELECT user2_id FROM block_list WHERE user1_id = ?
		UNION
		SELECT user1_id FROM block_list WHERE user2_id = ?
	`
	l := make(map[string]bool)
	rows, err := db.Query(query, id, id)
	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, nil
	}

	for rows.Next() {
		var blockedId string
		if err := rows.Scan(&blockedId); err != nil {
			return nil, err
		}
		l[blockedId] = true
	}
	return l, nil
}
