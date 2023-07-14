package main

import (
	"database/sql"
	"net/http"
	"problem1/configs"
	"problem1/handler"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "minimal_sns!")
	})

	e.GET("/users/:id", handler.GetUserHandler(db))

	e.GET("/users", handler.GetUserListHandler(db))

	e.GET("/users/:id/friends", handler.GetFriendListHandler(db))

	// クエリパラメータ limit, offset, hop を受け取ることができる
	e.GET("/users/:id/friends-of-friends", handler.GetFriendOfFriendListHandler(db))

	// e.GET("/get_friend_of_friend_list_paging", handler.GetFriendOfFriendListPagingHandler(db))

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
