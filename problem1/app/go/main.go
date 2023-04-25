package main

import (
	"database/sql"
	"fmt"
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
		fmt.Printf("%d\n", id)
		// id を使って，DB から友達リストを取得する
		friendList := []string{"friend1", "friend2", "friend3"}
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
