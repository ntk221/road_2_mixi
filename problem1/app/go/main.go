package main

import (
	"database/sql"
<<<<<<< HEAD
	"handler"
	"net/http"
	"problem1/configs"
	"problem1/handler"
=======
	"net/http"
	"problem1/configs"
	"problem1/dao"
>>>>>>> origin/main
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

	e.GET("/get_friend_list", handler.GetFriendListHandler(db))

	e.GET("/get_friend_of_friend_list", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.GET("/get_friend_of_friend_list_paging", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
