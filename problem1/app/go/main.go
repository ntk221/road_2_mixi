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

	e.GET("/get_friend_list", handler.GetFriendListHandler(db))

	e.GET("/get_friend_of_friend_list", handler.GetFriendOfFriendListHandler(db))

	e.GET("/get_friend_of_friend_list_paging", handler.GetFriendOfFriendListPagingHandler(db))

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
